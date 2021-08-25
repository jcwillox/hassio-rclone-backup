# Rclone Backup
Backup your Home Assistant configuration or backups to over 40 cloud providers using [Rclone](https://rclone.org/).

This pairs well with the custom integration [Auto Backup](https://github.com/jcwillox/hass-auto-backup) which provides a highly configurable way to create backups and have them deleted after a given period.

Rclone Backup can sync specific backups, e.g. backups starting with `AutoBackup*` to a cloud provider, and when that backup is deleted from Home Assistant it will be removed from the cloud provider as well.

You can also directly sync your Home Configuration e.g. `/config`, `/share`, `/ssl` to a cloud service or to another machine using SFTP. Rclone is smart and will only upload changed files.

# Installation

Add the repository URL under **Supervisor** → **Add-on store** → **⋮** → **Manage add-on repositories**

```
https://github.com/jcwillox/hassio-rclone-backup
```
