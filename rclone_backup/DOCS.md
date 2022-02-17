# Rclone Backup

Backup your Home Assistant configuration or backups to over 40 cloud providers using [Rclone](https://rclone.org/).

```yaml
jobs:
  - name: Sync Daily Backups
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

## Configuration

**Option:** `jobs`

A list of jobs to schedule, see the [job config](#job-config) section for a list of available options.

**Option:** `dry_run`

Trial run with no permanent changes, see what rclone would do without actually doing it.

**Option:** `run_once`

Will run all jobs with no schedule defined immediately then exit. Designed for use with the `hassio.addon_restart` service to trigger jobs.

**Option:** `config_path`

The location of the rclone config file, must be stored under `/config`.

**Option:** `rclone_config`

This allows you to specify the rclone config from within the addon configuration, defining it will take precedence over the `config_path` option.

```yaml
jobs:
  - ...
rclone_config: |
  [google]
  type = drive
  scope = drive
  token = REDACTED
  use_trash = false
```

**Option:** `no_rename`

Disable the renaming of backups before uploading them, with this enabled backups will use their slug id (their name on disk) instead of their friendly name.

**Option:** `no_unrename`

Prevents renaming backups back to their original name, with this enabled the backup files will be renamed to use their friendly name and this will persist on disk.

*Note: this option can apparently cause issues with restoring backups but won't affect the actual integrity of the snapshots, exercise caution when using it.*

## Job Config

**Option:** `name`

Optionally you can provide a friendly name for the job, this can be useful to identify which job is being run when you have multiple.

**Option:** `schedule`

Specify when the rclone backup should run using cron syntax. If the `schedule` option is empty or undefined the job will be run when the addon starts.

**Option:** `command`

The rclone command to run e.g. `sync`, `copy`, `move`.

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

---

### Generating the rclone config

I'd recommend consulting the [rclone docs](https://rclone.org/docs/), but simply you will need to install rclone locally and run the `rclone config` command, then follow the prompts to configure the services you need. The generated config is located by default at `$HOME/.config/rclone/rclone.conf` (alternatively run `rclone config file` and it will print the config location).
