import json
import os
import subprocess
import tarfile
from datetime import datetime
from glob import glob
from json import JSONDecodeError
from os import path, chdir
from os.path import isdir, splitext
from os.path import isfile
from subprocess import CalledProcessError
from typing import Dict

from slugify import slugify

INSTALL_PATH = "/opt/rclone-backup"
CONFIG_PATH = "/data/options.json"
BACKUP_PATH = "/backup"
ALLOWED_COMMAND = ["sync", "copy"]
ALLOWED_SOURCE_PATHS = ("/backup", "/config", "/share", "/ssl", "/media")


def now():
    return datetime.now().strftime("%Y-%m-%d %H:%M:%S")


def main():
    with open(CONFIG_PATH) as file:
        config = json.loads(file.read())

    start_time = now()

    command = config["command"]
    sources = config["sources"]
    destination = config["destination"]
    rclone_config_path = config["config_path"]

    if command not in ALLOWED_COMMAND:
        print(
            f"[rclone-backup] Given command is not allowed! Allowed commands: {ALLOWED_COMMAND}"
        )
        exit(1)

    if not isfile(rclone_config_path):
        print(
            f"[rclone-backup] Given rclone config file '{rclone_config_path}' does not exist!"
        )
        exit(1)

    with open(rclone_config_path) as file:
        if not any(line.startswith(f"[{destination.split(':')[0]}]") for line in file):
            print(
                f"[rclone-backup] Did not find any rclone configuration matching '{destination}'!"
            )
            exit(1)

    for source in sources:
        source = path.join("/", source)
        if not source.startswith(ALLOWED_SOURCE_PATHS):
            print(
                f"[rclone-backup] Given source '{source}' is not allowed! Allowed sources: {ALLOWED_SOURCE_PATHS}"
            )
            continue
        elif not isdir(source):
            print(f"[rclone-backup] Given source '{source}' directory does not exist!")
            continue

        subfolder = ""
        if len(sources) > 1:
            subfolder = f"{source}"

        renamed_backups = {}
        if source.startswith(BACKUP_PATH):
            renamed_backups = rename_backups()

        cmd = [
            "rclone",
            command,
            source,
            destination + subfolder,
            "--config",
            rclone_config_path,
            "--verbose",
        ]

        for include in config["include"]:
            cmd.append("--include")
            cmd.append(include)

        for exclude in config["exclude"]:
            cmd.append("--exclude")
            cmd.append(exclude)

        if config.get("dry_run"):
            cmd.append("--dry-run")

        for flag in config["flags"]:
            cmd.append(flag)

        try:
            subprocess.run(cmd, stdout=True, stderr=True, check=True)
        except CalledProcessError:
            print(f"[rclone-backup] Rclone failed!")

        if renamed_backups:
            undo_rename_backups(renamed_backups)

    print(f"[rclone-backup] Started at {start_time}")
    print(f"[rclone-backup] Finished at {now()}")
    print("[rclone-backup] Done!")


def get_backup_info(filename) -> (str, str):
    with tarfile.open(filename, "r:") as file:
        backup_config = "./backup.json"
        try:
            file.getmember(backup_config)
        except KeyError:
            backup_config = "./snapshot.json"
        data = json.loads(file.extractfile(backup_config).read())
    return data["name"], data["slug"]


def rename_backups() -> Dict[str, str]:
    print(f"[rclone-backup-rename] Starting at {now()}")
    renamed_backups: Dict[str, str] = {}
    chdir(BACKUP_PATH)

    for filename in glob("*.tar"):
        try:
            name, slug = get_backup_info(filename)
        except (tarfile.TarError, JSONDecodeError):
            continue
        friendly_filename = slugify(name, lowercase=False, separator="_") + ".tar"
        # we only want to rename backups that are named with their slug
        if splitext(filename)[0] == slug and not isfile(friendly_filename):
            # track renamed backups to restore their names later
            renamed_backups[filename] = friendly_filename
            try:
                os.rename(filename, friendly_filename)
            except OSError:
                print(
                    f"[rclone-backup-rename] Failed to rename '{filename}' to '{friendly_filename}'"
                )

    print(f"[rclone-backup-rename] Renamed {len(renamed_backups)} backups")
    print(f"[rclone-backup-rename] Finished at {now()}")
    return renamed_backups


def undo_rename_backups(renamed_backups: Dict[str, str]):
    print(f"[rclone-backup-undo-rename] Starting at {now()}")
    chdir(BACKUP_PATH)
    for name_slug, name_friendly in renamed_backups.items():
        try:
            os.rename(name_friendly, name_slug)
        except OSError:
            print(
                f"[rclone-backup-undo-rename] Failed to rename '{name_friendly}' to '{name_slug}'"
            )
    print(f"[rclone-backup-undo-rename] Finished at {now()}")


if __name__ == "__main__":
    main()
