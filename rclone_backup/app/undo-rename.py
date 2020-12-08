#!/usr/bin/env python3
import json
import tarfile
import os

from os import listdir
from os.path import isfile


BACKUP_PATH = "/backup"

os.chdir(BACKUP_PATH)

for snapshot in listdir():
    with tarfile.open(snapshot, "r:") as file:
        data = json.loads(file.extractfile("./snapshot.json").read())
    filename = data["slug"] + ".tar"
    if snapshot != filename and not isfile(filename):
        os.rename(snapshot, filename)
        print(f"[RENAMED] '{snapshot}' to '{filename}'")