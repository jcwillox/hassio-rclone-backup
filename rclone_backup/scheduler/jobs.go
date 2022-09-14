package main

import (
	"github.com/jcwillox/emerald"
	"os"
	"os/exec"
	"strings"
	"time"
)

// CreateJob create run job closure with the job config
func CreateJob(job JobConfig) func() {
	return func() {
		if len(job.Sources) > 1 && len(job.Destinations) > 1 {
			// multiple destinations and multiple sources
			for _, destination := range job.Destinations {
				for _, source := range job.Sources {
					RunJob(job, source, destination+source)
				}
			}
		} else if len(job.Sources) > 1 {
			// multiple sources
			// multiple sources to single destination
			if len(job.Destinations) > 0 {
				job.Destination = job.Destinations[0]
			}
			for _, source := range job.Sources {
				RunJob(job, source, job.Destination+source)
			}
			return
		} else if len(job.Destinations) > 1 {
			// multiple destinations
			// multiple destinations to single source
			if len(job.Sources) > 0 {
				job.Source = job.Sources[0]
			}
			for _, destination := range job.Destinations {
				RunJob(job, job.Source, destination)
			}
		} else {
			// single source
			// single destination
			// single source to single destination
			if len(job.Destinations) > 0 {
				job.Destination = job.Destinations[0]
			}
			if len(job.Sources) > 0 {
				job.Source = job.Sources[0]
			}
			RunJob(job, job.Source, job.Destination)
		}
	}
}

func RunJob(job JobConfig, source string, destination string) {
	// generate rclone command
	args := []string{job.Command, source}

	// destination is not required
	if destination != "" {
		args = append(args, destination)
	}

	args = append(args, "--verbose", "--config", config.ConfigPath)

	for _, inclusion := range job.Include {
		args = append(args, "--include", inclusion)
	}

	for _, exclusion := range job.Exclude {
		args = append(args, "--exclude", exclusion)
	}

	if config.DryRun {
		args = append(args, "--dry-run")
	}

	// append any extra flags
	args = append(args, FlagMapToList(config.Flags)...)
	args = append(args, config.ExtraFlags...)
	args = append(args, FlagMapToList(job.Flags)...)
	args = append(args, job.ExtraFlags...)

	Infoln("running", JobInfo(job, "job", source, destination))
	Debugln("rclone", args)

	start := time.Now()

	var undoRename func()
	if strings.HasPrefix(source, BackupPath) && !config.NoRename {
		var err error
		undoRename, err = RenameBackups(config.NoSlugify)
		if err != nil {
			Errorln("failed to rename backups, aborting upload", err)
			return
		}
	}

	emerald.Print(emerald.Blue)

	cmd := exec.Command("rclone", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		Errorln("failed to run rclone command", err)
		return
	}

	emerald.Print(emerald.Reset)

	if undoRename != nil && !config.NoUnrename {
		undoRename()
	}

	Infoln("finished in", boldCyan(FormatDuration(time.Since(start))))
}
