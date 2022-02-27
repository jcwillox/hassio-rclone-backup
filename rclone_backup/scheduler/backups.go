package main

import (
	"archive/tar"
	"encoding/json"
	"github.com/gosimple/slug"
	"github.com/jcwillox/emerald"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func RenameBackups() (func(), error) {
	renamed := make(map[string]string)
	files, err := filepath.Glob(filepath.Join(BackupPath, "*.tar"))
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		config, err := GetBackupConfig(file)
		if err != nil {
			return nil, err
		}
		if config == nil {
			continue
		}

		friendlyName := slug.Make(config.Name) + ".tar"

		// we only want to rename backups that are named with their slug
		fileName := strings.TrimSuffix(filepath.Base(file), ".tar")
		dest := filepath.Join(BackupPath, friendlyName)
		if fileName == config.Slug {
			if stat, _ := os.Stat(dest); stat == nil {
				err := os.Rename(file, dest)
				if err != nil {
					Errorln("failed to rename backup", emerald.HighlightPathStat(file), Arrow, emerald.HighlightPathStat(dest))
					return nil, err
				}
				renamed[file] = dest
			}
		}
	}

	Infoln("renamed", boldCyan(strconv.Itoa(len(renamed))), emerald.Green+"backups")

	return func() {
		for file, dest := range renamed {
			err := os.Rename(dest, file)
			if err != nil {
				Errorln("failed to unrename backup", emerald.HighlightPathStat(dest), Arrow, emerald.HighlightPathStat(file))
			}
		}
		Infoln("unrenamed", boldCyan(strconv.Itoa(len(renamed))), emerald.Green+"backups")
	}, err
}

func GetBackupConfig(file string) (*BackupConfig, error) {
	reader, err := os.Open(file)
	defer reader.Close()
	if err != nil {
		return nil, err
	}
	tr := tar.NewReader(reader)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		if header.Name == "./backup.json" || header.Name == "./snapshot.json" {
			data, err := io.ReadAll(tr)
			if err != nil {
				return nil, err
			}
			config := &BackupConfig{}
			err = json.Unmarshal(data, config)
			return config, err
		}
	}
}
