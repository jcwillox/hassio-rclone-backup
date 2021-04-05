#!/usr/bin/with-contenv bashio

echo $(rclone -V)

echo "$(bashio::config schedule) stdout python '/run-rclone.py'" > /etc/crontabs/root

cat /etc/crontabs/root

exec /sbin/tini -s -- /usr/sbin/crond -f
