# Rclone Backup

Backup your Home Assistant configuration or snapshots to over 40 cloud providers using [Rclone](https://rclone.org/).

## Configuration

```yaml
rclone:
  enabled: true
  schedule: 10 4 * * *
  command: sync
  source: /backup
  destination: 'google:/Backup/Home Assistant'
  flags: ''
  include:
    - DailyBackup*
  exclude: []
  dry_run: false
  config_path: /share/rclone.conf
rename:
  enabled: true
  schedule: 0 * * * *
```

---

### `rclone`

**Option:** `schedule`

Specify when the rclone backup should run using cron syntax.

**Option:** `command`

The rclone command to run e.g. `sync` or `copy`.

**Option:** `source`

The directory to read from

**Option:** `destination`

The location to write to in the format `remote:path`, see [rclone docs](https://rclone.org/docs).

**Option:** `flags`

Extra flags to give to the rclone command, see [rclone flags](https://rclone.org/flags).

**Option:** `include`

List of files or folders to include, see [rclone filtering](https://rclone.org/filtering).

**Option:** `exclude`

List of files or folders to exclude, see [rclone filtering](https://rclone.org/filtering).

**Option:** `dry_run`

Trial run with no permanent changes, see what rclone would do without actually doing it.

**Option:** `config_path`

The location of the rclone config file.

---

### `rename`

Renames snapshots in /backup from `slug.tar` e.g. `dc7d0645.tar` to use their name e.g. `DailyBackup_Monday.tar`. This is necessary for the `include`/`exclude` options to work correctly.

**Option:** `schedule`

Specify when the rename action should run using cron syntax.

Note: *this will always be run before the rclone backup when enabled.*