import json
import os
import tarfile
from datetime import datetime
from os import listdir
from os.path import isfile

CONFIG_PATH = "/data/options.json"
BACKUP_PATH = "/backup"

with open(CONFIG_PATH) as file:
    config = json.loads(file.read())

if not config["rename"]["enabled"]:
    exit(0)

print(f"[UNDO-RENAME] Running {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")

os.chdir(BACKUP_PATH)

for snapshot in listdir():
    with tarfile.open(snapshot, "r:") as file:
        data = json.loads(file.extractfile("./snapshot.json").read())
    filename = data["slug"] + ".tar"
    if snapshot != filename and not isfile(filename):
        os.rename(snapshot, filename)
        print(f"[UNDO-RENAME] '{snapshot}' to '{filename}'")
