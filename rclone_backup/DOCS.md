# Rclone Backup

Backup your Home Assistant configuration or backups to over 40 cloud providers using [Rclone](https://rclone.org/).

## Configuration

```yaml
schedule: 10 4 * * *
command: sync
sources:
  - /backup
destination: 'google:/Backup/Home Assistant'
include:
  - DailyBackup*
exclude: []
flags: []
dry_run: false
config_path: /config/rclone.conf
```

---

**Option:** `schedule`

Specify when the rclone backup should run using cron syntax.

**Option:** `command`

The rclone command to run e.g. `sync` or `copy`.

**Option:** `sources`

List of directories to read from must one of or a subdirectory of `/backup`, `/config`, `/share`, `/ssl`, `/media`.

*When specifying multiple sources the files will be stored under `destination/source` otherwise they will be directly under `destination`.*

**Option:** `destination`

The location to write to in the format `remote:path`, see [rclone docs](https://rclone.org/docs).

**Option:** `include`

List of files or folders to include, see [rclone filtering](https://rclone.org/filtering).

**Option:** `exclude`

List of files or folders to exclude, see [rclone filtering](https://rclone.org/filtering).

**Option:** `flags`

List of extra flags to give to the rclone command, see [rclone flags](https://rclone.org/flags).

*For example, you may want to add `--drive-use-trash=false` when using google drive so rclone immediately deletes files instead of sending them to the trash.*

**Option:** `dry_run`

Trial run with no permanent changes, see what rclone would do without actually doing it.

**Option:** `config_path`

The location of the rclone config file, must be stored under `/config`.
