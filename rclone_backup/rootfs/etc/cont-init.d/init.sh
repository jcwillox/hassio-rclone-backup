#!/usr/bin/with-contenv bashio

CONFIG_PATH="/config/rclone.conf"

if bashio::config.has_value "rclone_config"; then
  # write rclone config from addon config
  CONFIG_PATH="/root/.config/rclone/rclone.conf"
  echo -n "$(bashio::config "rclone_config")" > "$CONFIG_PATH"
elif bashio::config.has_value "config_path"; then
  CONFIG_PATH=$(bashio::config "config_path")
fi

echo -n "$CONFIG_PATH" > /var/run/s6/container_environment/RCLONE_CONFIG
