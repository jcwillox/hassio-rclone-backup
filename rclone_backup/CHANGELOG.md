## What’s changed

This is a major release as it includes support for the Rclone Web UI using ingress. This means you can now configure your remotes through an easy-to-use interface built into Home Assistant. Additionally, all dependencies have been updated to their latest versions.

## ⚡ Features
* Support Rclone WebUI (c5f8790) (closes #3)
* Send status events to Home Assistant (f75dd5e) (closes #16)
	* This also includes a blueprint to handle sending notifications ([see more](https://github.com/jcwillox/hassio-rclone-backup/blob/main/rclone_backup/DOCS.md#events)).
	* This can be disabled with the `no_events` option.
* Improve slugification (1737c8e)
	* Disallowed characters are now replaced with underscores and multiple underscores in a row are removed.
	* This is a minor **breaking-change** as it is possible that some backups will now be named slightly differently, causing them to be re-uploaded, etc, as rclone will believe they are different files.
* Add `no_slugify` option (1708787) (fixes #14)
	* This shouldn't cause issues when targeting Linux-based systems, but certain characters in filenames such as `:` can cause issues with Windows systems.
* Include providers list in README.md (0a61393) (closes #12)
* Make `config_path` optional (e8a893e) (closes #21)

## Changes
* Always print list of jobs (935abfd)
* Add issue templates (25058a4)
* Update link in README.md (f7896f4)
* Bump dependencies (d8e706b)

**Full Changelog**: https://github.com/jcwillox/hassio-rclone-backup/compare/2.0.1...3.0.0
