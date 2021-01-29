import json
import subprocess
from datetime import datetime
from subprocess import STDOUT, CalledProcessError

CONFIG_PATH = "/data/options.json"

with open(CONFIG_PATH) as file:
    config = json.loads(file.read())

if not config["rclone"]["enabled"]:
    exit(0)

print(f"[RCLONE] Running {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")

print("[RCLONE] Renaming snapshots...")

try:
    output = subprocess.check_output("python /run-rename.py", stderr=STDOUT, shell=True)
    print(output.decode())
except CalledProcessError as ex:
    print(ex.output.decode())
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
        output = subprocess.check_output(cmd, stderr=STDOUT, shell=True)
        print(output.decode())
    except CalledProcessError as ex:
        print(ex.output.decode())
        print(f"[RCLONE] Rclone Failed!")


print("[RCLONE] Done!")
