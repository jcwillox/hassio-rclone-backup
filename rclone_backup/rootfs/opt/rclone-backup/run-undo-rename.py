import json
import tarfile
from datetime import datetime
from os import chdir
from os import listdir
from os import rename
from os.path import isfile

BACKUP_PATH = "/backup"

print(f"[rclone-backup-undo-rename] Starting at {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")

chdir(BACKUP_PATH)

for backup in listdir():
    with tarfile.open(backup, "r:") as file:
        backup_config = "./backup.json"
        try:
            file.getmember(backup_config)
        except KeyError:
            backup_config = "./snapshot.json"
        data = json.loads(file.extractfile(backup_config).read())
    filename = data["slug"] + ".tar"
    if backup != filename and not isfile(filename):
        rename(backup, filename)
        print(f"[rclone-backup-undo-rename] Renamed '{backup}' to '{filename}'")

print(f"[rclone-backup-undo-rename] Finished at {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
