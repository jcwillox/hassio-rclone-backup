#!/usr/bin/with-contenv bashio

bashio::log.info "Installed rclone version:"
bashio::log.info "$(rclone -V)"

bashio::log.info "starting scheduler..."
exec /usr/bin/scheduler
