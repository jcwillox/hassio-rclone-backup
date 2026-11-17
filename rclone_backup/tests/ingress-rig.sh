#!/usr/bin/env bash
# ==============================================================================
# Ingress rig for the add-on's nginx proxy
#
# Builds an HA-shaped Docker network (172.30.32.0/23) with the Supervisor's ingress proxy at
# 172.30.32.2 and another add-on at 172.30.33.10. rclone runs `rcd --rc-no-auth` on
# 127.0.0.1:5572 at the version pinned in the Dockerfile, like the add-on but without the Web UI
# files (`--rc-files`), which do not change how rc calls are routed. nginx from the add-on's base
# image shares its network namespace and loads the given nginx.conf.
#
# Usage: rclone_backup/tests/ingress-rig.sh [path/to/nginx.conf]
# Needs Docker. Prints one line per check and exits 1 if any check fails.
# ==============================================================================
set -euo pipefail

here=$(cd "$(dirname "$0")" && pwd)
conf=${1:-$here/../rootfs/etc/nginx/nginx.conf}
# Docker would create a missing bind-mount source as an empty directory
[ -f "$conf" ] || { echo "not a file: $conf" >&2 && exit 2; }
conf=$(cd "$(dirname "$conf")" && pwd)/$(basename "$conf")

case $(uname -m) in
    arm64 | aarch64) arch=aarch64 ;;
    x86_64 | amd64) arch=amd64 ;;
    *) echo "unsupported host architecture: $(uname -m)" >&2 && exit 2 ;;
esac
base=$(sed -n "s/^ *$arch: //p" "$here/../build.yaml")
rclone_version=$(sed -n 's/^ENV RCLONE_INSTALLED_VERSION=//p' "$here/../Dockerfile")

rig=rclone-rig-$$
addon_ip=172.30.33.2
supervisor_ip=172.30.32.2
other_addon_ip=172.30.33.10

cleanup() {
    docker rm -f "$rig-addon" "$rig-nginx" >/dev/null 2>&1 || true
    docker network rm "$rig" >/dev/null 2>&1 || true
}
trap cleanup EXIT

docker network create --subnet 172.30.32.0/23 --gateway 172.30.32.1 "$rig" >/dev/null
docker run -d --name "$rig-addon" --network "$rig" --ip "$addon_ip" \
    "rclone/rclone:$rclone_version" rcd --rc-no-auth --rc-addr 127.0.0.1:5572 >/dev/null
docker run -d --name "$rig-nginx" --network "container:$rig-addon" \
    -v "$conf:/rig/nginx.conf:ro" --entrypoint sh "$base" \
    -c 'apk add --no-cache nginx >/dev/null && exec nginx -c /rig/nginx.conf -g "daemon off;"' >/dev/null

# post <from ip> <path> [json body]: prints the HTTP status the proxy returns
post() {
    local body=${3-}
    [ -n "$body" ] || body='{}'
    docker run --rm --network "$rig" --ip "$1" --entrypoint curl "$base" \
        -s -o /dev/null -w '%{http_code}' --path-as-is --max-time 10 -X POST \
        -H 'Content-Type: application/json' --data "$body" "http://$addon_ip:8099$2" || true
}

failed=0
check() {
    if [ "$2" = "$3" ]; then
        echo "ok   $1 ($3)"
    else
        echo "FAIL $1: expected $2, got $3"
        failed=1
    fi
}

for _ in $(seq 1 60); do
    [ "$(post "$supervisor_ip" /rc/noopauth)" = 200 ] && break
    sleep 1
done

serve='{"type":"http","fs":"/tmp","addr":":0"}'
batch='{"inputs":[{"_path":"serve/start","type":"http","fs":"/tmp","addr":":0"}]}'

check "Supervisor reaches rclone through nginx" 200 "$(post "$supervisor_ip" /rc/noopauth)"
check "core/command stays open (Settings reads the config with cat)" 200 \
    "$(post "$supervisor_ip" /core/command '{"command":"version"}')"
check "another add-on is refused" 403 "$(post "$other_addon_ip" /rc/noopauth)"

# Every way to start a serve first, then count serves before anything could stop them
for path in /serve/start //serve/start /./serve/start /x/../serve/start /%73erve/start \
    /serve%2Fstart; do
    check "Supervisor blocked from $path" 403 "$(post "$supervisor_ip" "$path" "$serve")"
done
check "Supervisor blocked from job/batch (it can call serve/start)" 403 \
    "$(post "$supervisor_ip" /job/batch "$batch")"
check "rclone method names are case-sensitive, so /Serve/start is not a method" 404 \
    "$(post "$supervisor_ip" /Serve/start "$serve")"

# Asked inside the container, past nginx; a failed call must fail the check, not count as 0
if list=$(docker exec "$rig-addon" rclone rc serve/list); then
    serves=$(printf '%s' "$list" | grep -c '"id"' || true)
else
    serves="no answer from serve/list"
fi
check "no serve is running after the blocked calls" 0 "$serves"

for path in /serve/list /serve/types /serve/stop /serve/stopall; do
    check "Supervisor blocked from $path" 403 "$(post "$supervisor_ip" "$path")"
done

exit "$failed"
