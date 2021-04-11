#!/usr/bin/with-contenv bashio

CONFIG_PATH=$(bashio::config config_path)

# Ensure configuration exists
if ! bashio::fs.file_exists $CONFIG_PATH; then
    bashio::log.warning "Configuration file does not exist!"
    bashio::log.warning "If this is your first time starting this add-on ensure"
    bashio::log.warning "to create a valid rclone configuration in "$CONFIG_PATH" file"
    mkdir -p "$(dirname "$CONFIG_PATH")" && touch "$CONFIG_PATH" \
        || bashio::exit.nok "Failed to create rclone configuration"
else
    bashio::log.info "Rclone configuration found"
fi

bashio::log.info "Installed rclone version:"
bashio::log.info "$(rclone -V --config ${CONFIG_PATH})"

bashio::log.info "Setup rclone backup cron..."
echo "$(bashio::config schedule) python '/opt/rclone-backup/run-rclone.py'" > /etc/crontabs/root
