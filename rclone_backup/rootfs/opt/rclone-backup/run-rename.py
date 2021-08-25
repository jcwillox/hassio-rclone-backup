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

for backup in listdir():
    with tarfile.open(backup, "r:") as file:
        backup_config = "./backup.json"
        try:
            file.getmember(backup_config)
        except KeyError:
            backup_config = "./snapshot.json"
        data = json.loads(file.extractfile(backup_config).read())
    name, slug = data["name"], data["slug"]
    filename = slugify(name, lowercase=False, separator="_") + ".tar"
    if backup != filename and not isfile(filename):
        rename(backup, filename)
        print(f"[rclone-backup-rename] Renamed '{backup}' to '{filename}'")

print(f"[rclone-backup-rename] Finished at {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
