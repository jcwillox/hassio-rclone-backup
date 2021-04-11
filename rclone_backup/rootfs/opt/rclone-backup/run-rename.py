import json
import tarfile
from datetime import datetime
from os import chdir
from os import listdir
from os import rename
from os.path import isfile

from slugify import slugify

BACKUP_PATH = "/backup"

print(f"[rclone-backup-rename] Starting at {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")

chdir(BACKUP_PATH)

for snapshot in listdir():
    with tarfile.open(snapshot, "r:") as file:
        data = json.loads(file.extractfile("./snapshot.json").read())
    name, slug = data["name"], data["slug"]
    filename = slugify(name, lowercase=False, separator="_") + ".tar"
    if snapshot != filename and not isfile(filename):
        rename(snapshot, filename)
        print(f"[rclone-backup-rename] Renamed '{snapshot}' to '{filename}'")

print(f"[rclone-backup-rename] Finished at {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
