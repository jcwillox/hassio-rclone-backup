#!/usr/bin/with-contenv bashio

echo $(rclone -V)

echo "$(bashio::config rename.schedule) tostdout python '/run-rename.py'" > /etc/crontabs/root
echo "$(bashio::config rclone.schedule) tostdout python '/run-rclone.py'" >> /etc/crontabs/root

cat /etc/crontabs/root

exec /sbin/tini -s -- /usr/sbin/crond -f