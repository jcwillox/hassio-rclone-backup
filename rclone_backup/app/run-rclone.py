import json
import subprocess
import sys
from datetime import datetime
from subprocess import CalledProcessError

CONFIG_PATH = "/data/options.json"
BACKUP_PATH = "/backup"

with open(CONFIG_PATH) as file:
    config = json.loads(file.read())

if not config["rclone"]["enabled"]:
    exit(0)

rclone = config["rclone"]

if BACKUP_PATH in rclone["sources"]:    
    try:
        subprocess.run(
            [sys.executable, "/run-rename.py"], stdout=True, stderr=True, check=True
        )
    except CalledProcessError as ex:
        print(f"[RCLONE] Rename failed!")

print(f"[RCLONE] Running {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")

if rclone["enabled"]:
    command = rclone["command"]
    sources = rclone["sources"]
    destination = rclone["destination"]
    config_path = rclone["config_path"]

    for source in rclone["sources"]:
        subfolder = ""
        if len(rclone["sources"]) > 1:
            subfolder = f"{source}"

        cmd = f"rclone {command} '{source}' '{destination}{subfolder}' --config '{config_path}' --verbose"

        for include in rclone["include"]:
            cmd += f" --include='{include}'"

        for exclude in rclone["exclude"]:
            cmd += f" --exclude='{exclude}'"

        if rclone.get("dry_run"):
            cmd += " --dry-run"

        if rclone.get("flags"):
            cmd += " " + rclone["flags"]

        print(f"[RCLONE] {cmd}")

        try:
            subprocess.run(cmd, stdout=True, stderr=True, check=True, shell=True)
        except CalledProcessError as ex:
            print(f"[RCLONE] Rclone failed!")

print("[RCLONE] Done!")

if BACKUP_PATH in rclone["sources"]:
    try:
        subprocess.run(
            [sys.executable, "/run-undo-rename.py"], stdout=True, stderr=True, check=True
        )
    except CalledProcessError as ex:
        print(f"[RCLONE] Undo rename failed!")
