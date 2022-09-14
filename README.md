# Rclone Backup
Backup your Home Assistant configuration or backups to over [40 cloud providers](https://rclone.org/#providers) using [Rclone](https://rclone.org/).

This pairs well with the custom integration [Auto Backup](https://github.com/jcwillox/hass-auto-backup) which provides a highly configurable way to create backups and have them deleted after a given period.

Rclone Backup can sync specific backups, e.g. backups starting with `AutoBackup*` to a cloud provider, and when that backup is deleted from Home Assistant it will be removed from the cloud provider as well.

You can also directly sync your Home Configuration e.g. `/config`, `/share`, `/ssl`, `/media` to a cloud service or to another machine using SFTP. Rclone is smart and will only upload changed files.

## Installation

[![Open your Home Assistant instance and show the add add-on repository dialog with a specific repository URL pre-filled.](https://my.home-assistant.io/badges/supervisor_add_addon_repository.svg)](https://my.home-assistant.io/redirect/supervisor_add_addon_repository/?repository_url=https%3A%2F%2Fgithub.com%2Fjcwillox%2Fhassio-rclone-backup)
[![Open your Home Assistant instance and show the dashboard of a Supervisor add-on.](https://my.home-assistant.io/badges/supervisor_addon.svg)](https://my.home-assistant.io/redirect/supervisor_addon/?repository_url=https%3A%2F%2Fgithub.com%2Fjcwillox%2Fhassio-rclone-backup&addon=19a172aa_rclone_backup)

Add the repository URL under **Supervisor** → **Add-on store** → **⋮** → **Manage add-on repositories**

```
https://github.com/jcwillox/hassio-rclone-backup
```

## Example

Rclone is a powerful tool, you could for example use the `crypt` and `googledrive` remotes to automatically encrypt your backups and upload them to google drive.

**`rclone.conf`**

```ini
[google]
type = drive
scope = drive
token = REDACTED
; immediately delete backups instead of sending them to the trash
use_trash = false

[hassbackup]
type = crypt
remote = google:Backup/Home Assistant
filename_encryption = off
directory_name_encryption = false
password = REDACTED
password2 = REDACTED
```

**Addon configuration**

```yaml
jobs:
  - name: Sync Daily Backups
    schedule: 10 4 * * *
    command: sync
    sources:
      - /backup
    destination: 'hassbackup:'
    include:
      - DailyBackup*
    exclude: []
    # we can also disable google drive trash using flags
    flags:
      drive-use-trash: false
dry_run: false
config_path: /config/rclone.conf
```

## Providers

This is a list of providers this addon supports synchronizing backups with, for an up-to-date list see [rclone.org/#providers](https://rclone.org/#providers).

1Fichier, Alibaba Cloud (Aliyun) Object Storage System (OSS), Amazon Drive (See note), Amazon S3, Backblaze B2, Box, Ceph, Citrix ShareFile, C14, DigitalOcean Spaces, Dreamhost, Dropbox, Enterprise File Fabric, FTP, Google Cloud Storage, Google Drive, Google Photos, HDFS, HTTP, Hubic, Jottacloud, IBM COS S3, Koofr, Mail.ru Cloud, Memset Memstore, Mega, Memory, Microsoft Azure Blob Storage, Microsoft OneDrive, Minio, Nextcloud, OVH, OpenDrive, OpenStack Swift, Oracle Cloud Storage, ownCloud, pCloud, premiumize.me, put.io, QingStor, Rackspace Cloud Files, rsync.net, Scaleway, Seafile, SeaweedFS, SFTP, Sia, StackPath, SugarSync, Tardigrade, Tencent Cloud Object Storage (COS), Uptobox, Wasabi, WebDAV, Yandex Disk, Zoho WorkDrive, The local filesystem.
