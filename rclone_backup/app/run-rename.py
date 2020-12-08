#!/usr/bin/env python3
import json
import tarfile
import os

from os import listdir
from os.path import isfile
from slugify import slugify
from datetime import datetime


CONFIG_PATH = "/data/options.json"
BACKUP_PATH = "/backup"

os.chdir(BACKUP_PATH)

with open(CONFIG_PATH) as file:
    config = json.loads(file.read())

if not config["rename"]["enabled"]:
    exit(0)

print(f"[RENAME] Running {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")


def slugify_filename(name):
    return slugify(name, lowercase=False, separator="_") + ".tar"


for snapshot in listdir():
    with tarfile.open(snapshot, "r:") as file:
        data = json.loads(file.extractfile("./snapshot.json").read())
    name, slug = data["name"], data["slug"]
    filename = slugify_filename(name)
    if snapshot != filename and not isfile(filename):
        os.rename(snapshot, filename)
        print(f"[RENAME] '{snapshot}' to '{filename}'")
