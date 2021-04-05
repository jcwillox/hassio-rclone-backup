import json
import subprocess
import sys
from datetime import datetime
from subprocess import CalledProcessError

CONFIG_PATH = "/data/options.json"

with open(CONFIG_PATH) as file:
    config = json.loads(file.read())

if not config["rclone"]["enabled"]:
    exit(0)

print(f"[RCLONE] Running {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")

print("[RCLONE] Renaming snapshots...")

try:
    subprocess.run(
        [sys.executable, "/run-rename.py"], stdout=True, stderr=True, check=True
    )
except CalledProcessError as ex:
    print(f"[RCLONE] Rename Failed!")

print("[RCLONE] Running rclone...")

with open(CONFIG_PATH) as file:
    config = json.loads(file.read())

rclone = config["rclone"]

if rclone["enabled"]:
    command = rclone["command"]
    source = rclone["source"]
    destination = rclone["destination"]
    config_path = rclone["config_path"]

    cmd = f"rclone {command} '{source}' '{destination}' --config '{config_path}' --verbose"

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
        print(f"[RCLONE] Rclone Failed!")

try:
    subprocess.run(
        [sys.executable, "/run-undo-rename.py"], stdout=True, stderr=True, check=True
    )
except CalledProcessError as ex:
    print(f"[RCLONE] Undo Rename Failed!")

print("[RCLONE] Done!")
