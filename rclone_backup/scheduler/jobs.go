package main

import (
	"fmt"
	"github.com/jcwillox/emerald"
	"os"
	"os/exec"
	"strings"
	"time"
)

// CreateJob create run job closure with the job config
func CreateJob(job JobConfig) func() {
	return func() {
		for _, run := range JobRuns(job) {
			RunJob(job, run[0], run[1])
		}
	}
}

// JobRuns the source and destination of each rclone run in a job, the destination is empty when
// the job has none (delete)
func JobRuns(job JobConfig) [][2]string {
	var runs [][2]string
	if len(job.Destinations) == 0 {
		// no destination, each source on its own
		for _, source := range job.Sources {
			runs = append(runs, [2]string{source, ""})
		}
	} else if len(job.Sources) > 1 {
		// multiple sources, each under destination/source
		for _, destination := range job.Destinations {
			for _, source := range job.Sources {
				runs = append(runs, [2]string{source, destination + source})
			}
		}
	} else {
		// single source to each destination
		for _, destination := range job.Destinations {
			runs = append(runs, [2]string{job.Sources[0], destination})
		}
	}
	return runs
}

// BuildArgs generate the rclone command arguments for a job
func BuildArgs(job JobConfig, source string, destination string) []string {
	args := []string{job.Command, source}

	// destination is not required
	if destination != "" {
		args = append(args, destination)
	}

	args = append(args, "--verbose")

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
	return args
}

func RunJob(job JobConfig, source string, destination string) {
	// generate rclone command
	args := BuildArgs(job, source, destination)

	Infoln("running", JobInfo(job, "job", source, destination))
	Debugln("rclone", args)

	start := time.Now()

	var undoRename func()
	if strings.HasPrefix(source, BackupPath) && !config.NoRename {
		var err error
		undoRename, err = RenameBackups(config.NoSlugify)
		if err != nil {
			msg := fmt.Sprintf("failed to rename backups, aborting upload: %s", err)
			Errorln(msg)
			FireJobEvent(EventJobFailed, job, source, destination, start, msg)
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
		msg := fmt.Sprintf("failed to run rclone command: %s", err)
		Errorln(msg)
		FireJobEvent(EventJobFailed, job, source, destination, start, msg)
		return
	}

	emerald.Print(emerald.Reset)

	if undoRename != nil && !config.NoUnrename {
		undoRename()
	}

	Infoln("finished in", boldCyan(FormatDuration(time.Since(start))))
	FireJobEvent(EventJobSuccessful, job, source, destination, start, "")
}
