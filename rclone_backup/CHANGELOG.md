## What’s changed

This major release has many significant changes, in particular, the entire core has been rewritten in golang, this was mainly due to the fact that Go actually has a good cron-syntax based scheduling library. The Go program now handles scheduling instead of cron, which is what has enabled support for multiple jobs. Otherwise, I've greatly improved logging, error handling, and made it so you can run pretty much any rclone command you'd like to.

## 🚨 Breaking Changes
* Support multiple scheduled jobs (closes #6)
  * Jobs are now specified as a list, and can optionally have a name. You will need to manually migrate your configuration.
   ```yaml
  # before
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
  ```yaml
  # after
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
      flags: {}
  dry_run: false
  config_path: /config/rclone.conf
  ```
* Flags are now specified as a map (d803e61)
  * Flags are now specified as key-value pairs, the current `flags` option has been renamed to `extra_flags` and a new `flags` option which expects a map has been introduced.
  ```yaml
  # before
  flags: 
    - --drive-use-trash=false
  ```
  ```yaml
  # after
  flags:
    drive-use-trash: false
  ```
* Renamed `disable_rename` and `disable_undo_rename` to `no_rename` and `no_unrename`.
  ```yaml
  # before
  disable_rename: false
  disable_undo_rename: false
  ```
  ```yaml
  # after
  no_rename: false
  no_unrename: false
  ```

## ⚡ Features
* Support running jobs on startup (closes #5)
* Allow multiple sources and multiple destinations (f4ed76e)
* Allow only the source to be specified
  * this means you can now use commands like `ls`, `purge`, `delete`.
* Allow sources to be remotes
* Added global `flags` option
* Added name option for jobs
* Added `run_once` option (closes #5)
  * this doesn't work exactly like you might expect at the moment as the program will exit but the addon will continue running, it will work well with the `hassio.addon_restart` service.
* Allow rclone config to be configured from the UI (1c505d7) (closes #13)

## Changes
* Make include/exclude optional (02f99de)
* Improve source/destination validation (a6c5b92)
* Migrate build and config files to YAML (98aba0c)
* Rewrite core in golang (70fdff5)

**Full Changelog**: https://github.com/jcwillox/hassio-rclone-backup/compare/1.2.0...2.0.0
