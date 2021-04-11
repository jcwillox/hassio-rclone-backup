#!/usr/bin/with-contenv bashio

config_path=$(bashio::config config_path)

bashio::log.info "Installed rclone version:"
bashio::log.info "$(rclone -V --config ${config_path})"

bashio::log.info "Setup rclone backup cron..."
echo "$(bashio::config schedule) python '/opt/rclone-backup/run-rclone.py'" > /etc/crontabs/root
