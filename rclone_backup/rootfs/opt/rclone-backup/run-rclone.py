import json
import subprocess
import sys
from datetime import datetime
from os import path
from os.path import isdir
from os.path import isfile
from subprocess import CalledProcessError

INSTALL_PATH = "/opt/rclone-backup"
CONFIG_PATH = "/data/options.json"
BACKUP_PATH = "/backup"
ALLOWED_COMMAND = ["sync", "copy"]
ALLOWED_SOURCE_PATHS = ("/backup", "/config", "/share", "/ssl")


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

        if source.startswith(BACKUP_PATH):
            try:
                subprocess.run(
                    [sys.executable, INSTALL_PATH + "/run-rename.py"],
                    stdout=True,
                    stderr=True,
                    check=True,
                )
            except CalledProcessError:
                print(f"[rclone-backup] Rename failed!")

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

        if source.startswith(BACKUP_PATH):
            try:
                subprocess.run(
                    [sys.executable, INSTALL_PATH + "/run-undo-rename.py"],
                    stdout=True,
                    stderr=True,
                    check=True,
                )
            except CalledProcessError:
                print(f"[rclone-backup] Undo rename failed!")

    print(f"[rclone-backup] Started at {start_time}")
    print(f"[rclone-backup] Finished at {now()}")
    print("[rclone-backup] Done!")


if __name__ == "__main__":
    main()
