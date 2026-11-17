# Rclone Backup

Backup your Home Assistant configuration or backups to over [70 cloud providers](https://rclone.org/#providers) using [Rclone](https://rclone.org/).

The app has two parts:

- **Jobs** that you write in the app configuration. They copy or sync files on a schedule.
- **The Web UI**, where you add remotes, browse files and see what rclone is doing.

If this is your first time, do the steps in **Before you start**. The app starts with no jobs, so you
can add a remote in the Web UI first, and then a job that uses it.

A typical configuration looks like this. It uses a remote called `google`, so add that remote before
you use it (see **Before you start**). If a job uses a remote that does not exist, the app stops when
it starts.

```yaml
jobs:
  - name: Sync Daily Backups
    schedule: 10 4 * * *
    command: sync
    sources:
      - /backup
    destination: 'google:Backup/Home Assistant'
    exclude: []
    flags: {}
dry_run: false
config_path: /homeassistant/rclone.conf
```

## Before you start

The app reads your remotes (your storage accounts) from a file called `rclone.conf`. You do not need
to create this file. The Web UI creates it when you add your first remote.

Every job names a remote. If a job names a remote that does not exist, the app stops when it starts,
and you cannot open the Web UI. So add your remote in the Web UI first, and add your job after that.

*Note: in version 3.5.0 and later, the app has no example job. The update does not change the jobs
that you saved in the **Configuration** tab.*

### 1. Start the app

1. Start the app. It has no jobs yet, so it only starts the Web UI.
2. Select **Open Web UI**. The **Dashboard** shows no remotes yet.

### 2. Add a remote in the Web UI

1. Select **Remotes** → **Add New Remote** and follow the steps in **Remotes** (under **Using the Web
   UI**). Give the remote a short name without spaces, for example `mydrive`.
2. If the provider logs you in through a page in your browser, for example Google Drive, OneDrive or
   Box, the Web UI asks you to run `rclone authorize` on a computer and paste the result. See
   **Providers with a browser login**.
3. Test the remote. In **Remotes**, select the remote. The Web UI lists its top-level folders. If you
   see an error, edit the remote and fix it before you continue.

The Web UI saves the remote in `/homeassistant/rclone.conf`. This is the default value of the
`config_path` option, so you do not need to change that option.

*Note: the file contains the passwords and tokens for your accounts. Keep it as safe as a password.*

### 3. Add your jobs and restart the app

1. Go to the app's **Configuration** tab.
2. Add a job that uses the name of your remote. For example:

   ```yaml
   jobs:
     - name: Sync Daily Backups
       schedule: 10 4 * * *
       command: sync
       sources:
         - /backup
       destination: 'mydrive:Backup/Home Assistant'
   ```

   The name must match the name of the remote exactly, and a colon (`:`) must follow it. The job has
   no filter, so it copies every backup in `/backup`, including the backups of single apps. To copy
   only some backups, see the `include` option.
3. Set `dry_run` to `false`. The default is `true`. With `true`, rclone only shows in the log what it
   would do, and it does not copy or delete files. You can keep `true` for the first run, check the
   log, and then change it to `false`.
4. Restart the app.
5. If the app stops, open its **Log** tab. The log names the job and the problem, for example a remote
   that it did not find. Correct the job in the **Configuration** tab and start the app again.

The app reads and checks the jobs only when it starts. After you change a job or add a remote that a
job uses, restart the app.

### If you already have an rclone.conf

If you already use rclone on a computer, you can use its remotes. Run `rclone config file` on the
computer to find the file. Then copy its contents to `/homeassistant/rclone.conf` with an app that can
edit files, for example the
[File editor](https://my.home-assistant.io/redirect/supervisor_addon/?addon=core_configurator) app.
Or paste the contents into the `rclone_config` option. A remote that you add in the Web UI is then
lost at the next restart, see **Remotes**.

## Using the Web UI

Select **Open Web UI** on the app's **Info** tab to open the
[rclone-web](https://github.com/rclone/rclone-web) interface inside Home Assistant.

![webui-example](https://github.com/jcwillox/hassio-rclone-backup/raw/main/example.png)

Not every tab works inside Home Assistant. This table shows what to expect:

| Tab           | Status        | What to know                                                           |
| ------------- | ------------- | ---------------------------------------------------------------------- |
| **Dashboard** | Works         | Shows remotes, mounts, transfer counters and the rclone version.       |
| **Explorer**  | Works, limits | You can only see the folders that the app's security profile allows.  |
| **Remotes**   | Works, limits | A provider with a browser login needs rclone on a computer too.        |
| **Mounts**    | Works, limits | Share on your network with another app. A restart removes it.          |
| **Transfers** | Works         | Copy and move files.                                                   |
| **Settings**  | Works         | Saving the config contents creates the file if it does not exist.      |

### Explorer

Explorer shows the files that rclone can see. The app runs rclone under an
[AppArmor](https://contributors.home-assistant.io/apps/guides/security/#apparmor) security profile, which
allows only these folders:

| Folder           | What it contains                                   | Access     |
| ---------------- | -------------------------------------------------- | ---------- |
| `/homeassistant` | Your Home Assistant configuration                  | Read/write |
| `/backup`        | Home Assistant backups                             | Read/write |
| `/share`         | The shared folder                                  | Read/write |
| `/media`         | The media folder                                   | Read/write |
| `/addons`        | Local apps                                         | Read/write |
| `/addon_configs` | The configuration folders of other apps            | Read/write |
| `/config`        | This app's own configuration folder                | Read/write |
| `/tmp`           | Temporary files                                    | Read/write |
| `/ssl`           | Certificates                                       | Read-only  |
| `/data`          | This app's options, including `rclone_config`      | Read-only  |

*Note: `/data/options.json` holds the app's options. If you set `rclone_config`, the file contains
the passwords and tokens for your accounts in plain text, and anyone who can open the Web UI can read
them in Explorer.*

Explorer shows other folders in the top level, for example `/etc` or `/root`, but you cannot open them.
When you try, Explorer shows an error and AppArmor blocks the access. The block is expected. To see
what AppArmor blocked, look in the logs for lines that contain `apparmor="DENIED"` and
`rclone_backup`.

### Remotes

A remote is a connection to a storage provider, for example Google Drive, an S3 bucket or an SFTP
server.

Some providers, for example S3, SFTP and WebDAV, use keys or passwords. You add them with the steps
below.

Other providers, for example Google Drive, OneDrive, Dropbox and Box, log you in through a page in
your browser. rclone runs inside Home Assistant, so it cannot open that page for you. The Web UI asks
you to log in on a computer that has rclone and a browser, then paste the result. See **Providers
with a browser login** below.

To add a remote in the Web UI:

1. Select **Remotes** → **Add New Remote**.
2. Give the remote a short name without spaces, for example `mydrive`. You use this name in your jobs.
3. Select the provider and fill in its settings. The [rclone docs](https://rclone.org/docs/) describe
   the settings for each provider.
4. Save the remote. The Web UI writes it to the file in `config_path`, by default
   `/homeassistant/rclone.conf`. If you set `rclone_config`, the app writes that option to the file
   again at each start, so a remote that you add in the Web UI is lost at the next restart. Add it to
   `rclone_config` instead.
5. If a job uses the new remote, restart the app.

#### Providers with a browser login

When you save a remote for such a provider, the Web UI shows **rclone needs more information**:

1. Read the command in the window, for example `rclone authorize "box"`.
2. On your computer, run that command. Use the same rclone version as the app if you can. A page opens
   in your browser. Log in and allow rclone to access your account.
3. The command prints a token between the lines `Paste the following into your remote machine --->`
   and `<---End paste`. The token starts with `{` and ends with `}`. Copy all of it.
4. Paste the token into the Web UI and select **Continue**.

If rclone has more questions, for example which OneDrive drive to use, the Web UI shows them one at a
time. Select **Cancel** to stop. The Web UI then removes the new remote.

When you edit one of these remotes, the Web UI asks if you want to replace the token. Select
`false: No` to keep your login. Select `true: Yes` only to log in again.

### Mounts

A mount makes a remote look like a local folder. Before you start, know these limits:

- **The mount stays inside this app.** It does not appear on your network by itself. To share it,
  mount it under `/share` or `/media`, then share that folder with another app. See **Share a
  mount on your network** below.
- **A mount does not survive a restart.** When the app restarts, or Home Assistant restarts, the mount
  is gone. You must create it again.
- **The mount folder must already exist.** The Web UI does not create it.

To create a mount:

1. Create an empty folder for the mount. Select **Explorer**, open `/share` (or `/media`), select
   **New Folder** and give it a name, for example `rclone`.
2. Select **Mounts** → **New Mount**.
3. Select the remote. In **Mount Point**, type the full path of the folder from step 1, for example
   `/share/rclone`.
4. Create the mount. The remote's files now show in that folder, inside this app.

After each restart of the app or of Home Assistant, do steps 2 to 4 again.

#### Share a mount on your network

Mount the remote under `/share` or `/media`, then install an app that shares that folder, for
example:

- **Samba share** (SMB), the official app.
- [NFS Server](https://github.com/sam-ward/nfs-ganesha-addon) (NFSv4), a community app.
- [FTP](https://github.com/hassio-addons/app-ftp) (FTP or FTPS), a community app.

When this app restarts, the mount is gone and the other app shows an empty folder until you create
the mount again.

### Settings

Settings shows the path and the contents of the rclone config file, and lets you edit the contents.
Saving the contents creates the file if it does not exist. If you set `rclone_config`, the app writes
that option to the file again at each start, so make the change in the option instead.

Settings also lets you change rclone options, for example a bandwidth limit. These options apply only
to the rclone that runs the Web UI, for example to a copy that you start in **Transfers**. They do not
apply to your jobs, because each job starts its own rclone command. A restart of the app also resets
them. To set an option for your jobs, use the `flags` or `extra_flags` option.

## Configuration

**Option:** `jobs`

A list of jobs to schedule, see the **Job Config** section for a list of available options.

**Option:** `flags`

Global map of flags to give the rclone command for every job, see [rclone flags](https://rclone.org/flags). Underscores are allowed for the flags and they will be replaced with dashes.

```yaml
jobs:
  - ...
flags:
  drive-use-trash: false
  min-age: 30d
```

*For example, you may want to add `drive-use-trash: false` when using google drive so rclone immediately deletes files instead of sending them to the trash.*

**Option:** `extra_flags`

List of flags to give the rclone command, applied globally to all jobs. For use when `flags` option isn't working, the list of flags are appended directly to the rclone command.

```yaml
jobs:
  - ...
extra_flags:
  - --drive-use-trash=false
  - --min-age=30d
```

**Option:** `dry_run`

Trial run with no permanent changes, see what rclone would do without actually doing it. The default
is `true`, so set it to `false` when your jobs are ready.

**Option:** `run_once`

Will run all jobs with no schedule defined immediately then exit. Designed for use with the `hassio.addon_restart` service to trigger jobs.

**Option:** `config_path`

The location of the rclone config file. The default is `/homeassistant/rclone.conf`. If the file does
not exist, the Web UI creates it when you add the first remote, see **Before you start**. Inside this app, `/config` is the app's own
configuration folder, not your Home Assistant configuration. If an older setup uses
`/config/rclone.conf`, change it to `/homeassistant/rclone.conf`.

**Option:** `rclone_config`

This allows you to specify the rclone config from within the addon configuration, defining it will take precedence over the `config_path` option.

The app writes this option to the rclone config file each time it starts. Changes that you make in
the Web UI are lost at the next restart, so make them in this option.

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

*Note: this option can apparently cause issues with restoring backups through the WebUI, however, selecting **⋮** → **Reload** should fix it. 
This won't affect the actual integrity of the backups.*

**Option:** `no_slugify`

This will disable the slugification of backup names. This means the **user** is responsible for ensuring their backup names are **valid** filenames for their filesystem and the destination filesystem.

**Option:** `no_events`

Disable sending completion and failure events to Home Assistant.

## Job Config

**Option:** `sources`

List of source locations in the format `remote:path` or for local paths just `path`, e.g. `/backup`, `/homeassistant`, `/share`, `/ssl`, `/media`.

**Option:** `source`

A slightly cleaner way to specify a single source, this option will override the `sources` option.

**Option:** `destinations`

A list of locations to write to in the format `remote:path` or for local paths just `path`, see [rclone docs](https://rclone.org/docs).

**Option:** `destination`

A slightly cleaner way to specify a single destination, this option will override the `destinations` option.

*Note: on SFTP, `remote:path` starts in your home folder on the server, and `remote:/path` starts at
the root of the server.*

> **Multiple Sources**
>
> When there are multiple sources each source will be backed-up to the destination under `destination/source`
> otherwise they will be directly under `destination`.
>
> **Multiple Destinations**
>
> All sources are backed-up to each destination following the rules for multiple sources mentioned above.

**Option:** `name`

Optionally you can provide a friendly name for the job, this can be useful to identify which job is being run when you have multiple.

**Option:** `schedule`

Specify when the rclone backup should run using cron syntax. If the `schedule` option is empty or undefined the job will be run when the addon starts.

| Form | Example | Meaning |
| --- | --- | --- |
| Five-field cron | `10 4 * * *` | 04:10 every day |
| `@hourly` | `"@hourly"` | At minute 0 of every hour |
| `@daily` or `@midnight` | `"@daily"` | 00:00 every day |
| `@weekly` | `"@weekly"` | 00:00 on Sunday |
| `@monthly` | `"@monthly"` | 00:00 on the 1st of the month |
| `@yearly` or `@annually` | `"@yearly"` | 00:00 on 1 January |
| `@every <duration>` | `"@every 6h"`, `"@every 1h30m"` | A fixed interval, counted from when the app starts, not aligned to the clock |
| Time zone prefix | `CRON_TZ=Europe/London 10 4 * * *` | Uses that time zone, not the app's local time |
| Empty or not set | | Runs once when the app starts |

The cron form has five fields (minute, hour, day of month, month, day of week). Seconds are not supported, so a six-field value fails.

*Note: in YAML, put quotes around a value that starts with `@` or `*`, for example `schedule: "@hourly"` or `schedule: "*/15 * * * *"`. Without quotes the YAML is not valid. An invalid schedule stops all jobs, so check the log after you change it.*

**Option:** `command`

The rclone command to run: `copy`, `sync`, `move` or `delete`. `copy`, `sync` and `move` go from the sources to the destination. `delete` acts on its sources only.

| Command | What it does |
| --- | --- |
| `copy` | Copies new and changed files to the destination. It never deletes anything there. |
| `sync` | Makes the destination match the source. Files at the destination that are not in the source are deleted, so backups that Home Assistant removes are removed from the destination too. |
| `move` | Copies to the destination, then deletes the files from the source. Do not use it on `/backup` unless you want no local backups. |
| `delete` | Deletes the files in the sources that match the job's filters, for example backups older than 14 days. |

The app checks each job when it starts. If a job is not correct, the app does not start, and the log names the job and the problem:

- `copy`, `sync` and `move` need a `destination` or `destinations`.
- `delete` must have no `destination` and no `destinations`.
- Each `delete` source must be a remote, for example `hassbackup:`. To keep fewer backups in Home Assistant, set the retention in Home Assistant itself.
- `delete` needs a filter, so that it does not delete every file. Use a `min-age` or `max-age` flag, `include` or `exclude`, or one of these filters in `extra_flags`.

Example: keep 14 days of backups on the remote. The `copy` job uploads backups each day at 04:10,
and the `delete` job removes the backups older than 14 days at 05:10.

```yaml
jobs:
  - name: Copy Backups
    schedule: 10 4 * * *
    command: copy
    sources:
      - /backup
    destination: 'hassbackup:'
    flags:
      max-age: 14d
  - name: Prune Old Backups
    schedule: 10 5 * * *
    command: delete
    sources:
      - 'hassbackup:'
    flags:
      min-age: 14d
```

Use the same age in both jobs. Without `max-age`, the `copy` job uploads again each backup that the
`delete` job removed, for as long as Home Assistant keeps it. Do not add a `delete` job for a remote
that a `sync` job writes to: `sync` already removes the backups that Home Assistant removes, and the
`delete` job removes backups that the next `sync` uploads again.

*Note: in version 3.5.0 and later, other commands such as `ls` are not allowed in a job. If a saved job uses one, the app does not start until you change it. To look at files once, use the Web UI.*

**Option:** `include`

List of files or folders to include, see [rclone filtering](https://rclone.org/filtering).

Without `include`, a job on `/backup` copies every backup: the full backups and the backups of single
apps. The names of the backups depend on the Home Assistant version and on who made the backup, and a
later version can change them. So look at the names before you add a filter (see the preview job
below), and check the log after a Home Assistant update.

*Note: these filters are matching the backups name on the filesystem, if you have not disabled renaming then this will be a slugified version of the backups name, this will effectively replace all special characters with an `_` and remove duplicate underscores.*

To preview the slugified names of your backups you can add a startup job like the following. It is a dry run to `/tmp`, so it copies nothing, and the log shows each name as `<name>: Skipped copy as --dry-run is set`.

```yaml
jobs:
  - command: copy
    name: Preview Backup Names
    sources:
      - /backup
    destination: /tmp
    flags:
      dry-run: true
```

**Option:** `exclude`

List of files or folders to exclude, see [rclone filtering](https://rclone.org/filtering).

**Option:** `flags`

Map of flags to give to the rclone command, see [rclone flags](https://rclone.org/flags).

**Option:** `extra_flags`

List of flags to give this job's rclone command only. For use when `flags` option isn't working, the list of flags are appended directly to the rclone command, after the global `flags` and `extra_flags`.

---

### Configuring Rclone Remotes

To add a remote in the Web UI, see **Remotes** (under **Using the Web UI**). For providers that need a browser login, the
Web UI asks you to paste a token from `rclone authorize` on your computer. The
[rclone docs](https://rclone.org/docs/) describe each provider.

---

## Events

> Blueprint to send notifications based on these events. [[source]](https://github.com/jcwillox/home-assistant-blueprints)
>
> [![Open your Home Assistant instance and show the blueprint import dialog with a specific blueprint pre-filled.](https://my.home-assistant.io/badges/blueprint_import.svg)](https://my.home-assistant.io/redirect/blueprint_import/?blueprint_url=https%3A//raw.githubusercontent.com/jcwillox/home-assistant-blueprints/main/automation/notify_rclone_backup.yaml)

The addon will fire events using Home Assistant's event bus whenever a job is completed or fails.

**Event:** `rclone_backup.job_successful`

**Event:** `rclone_backup.job_failed`

Both events will have the following attributes.

| Attribute     | Description                                            |
| ------------- | ------------------------------------------------------ |
| `name`        | The name of the job.                                   |
| `command`     | The rclone command that was run.                       |
| `source`      | The source location.                                   |
| `destination` | The destination location. (optional)                   |
| `error`       | The error message if the job failed. (optional)        |
| `duration`    | The duration of the job as a human string, eg. `1m2s`. |
| `seconds`     | The duration of the job in seconds.                    |

---

## Troubleshooting

### OneDrive: `Personal Vault` fails with `ObjectHandle is Invalid`

A job with a OneDrive Personal source fails, and the log shows lines like these:

```text
ERROR : Personal Vault: error reading source directory: couldn't list files: invalidRequest: invalidResourceId: ObjectHandle is Invalid
NOTICE: Failed to copy: couldn't list files: invalidRequest: invalidResourceId: ObjectHandle is Invalid
```

The other folders are copied, but the job ends as failed and fires `rclone_backup.job_failed`.

**Cause:** OneDrive Personal puts a locked folder named `Personal Vault` in the root of each account. While the vault is locked, OneDrive does not let rclone list it, and rclone cannot unlock it. Retries do not help. This is not a login problem: other folders on the same remote list correctly.

**Fix:** exclude the vault from the job.

```yaml
jobs:
  - name: Copy OneDrive to Drime
    command: copy
    source: 'onedrive:'
    destination: drime:Backup/OneDrive
    exclude:
      - "/Personal Vault/**"
```

The leading `/` matches only the vault in the root of the source, not other folders with the same name. You can keep the space in the name, because the app gives the value to rclone as written.

*Note: files in the Personal Vault are not backed up. To back them up, move them out of the vault in OneDrive.*

To test the filter on your computer before you change the job:

```sh
rclone copy onedrive: /path/to/test --exclude "/Personal Vault/**" --dry-run -v
```
