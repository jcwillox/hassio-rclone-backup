package main

import (
	"bufio"
	"github.com/jcwillox/emerald"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func PrintJobs(jobs []JobConfig) {
	lSchedule := 8
	lCommand := 0
	for _, job := range jobs {
		if len(job.Schedule) > lSchedule {
			lSchedule = len(job.Schedule)
		}
		if len(job.Command) > lCommand {
			lCommand = len(job.Command)
		}
	}

	Infoln("scheduled jobs:")

	for _, job := range config.Jobs {
		if job.Schedule == "" {
			job.Schedule = "@startup"
		}
		emerald.Print(job.Schedule, strings.Repeat(" ", lSchedule-len(job.Schedule)), " ")
		emerald.Print(emerald.Yellow, job.Command, emerald.Reset, strings.Repeat(" ", lCommand-len(job.Command)), " ")
		emerald.Println(JobInfo(job, ""))
	}
}

func FormatDuration(d time.Duration) string {
	scale := 100 * time.Second
	for scale > d {
		scale = scale / 10
	}
	return d.Round(scale / 100).String()
}

func ArrayContains(arr []string, s string) bool {
	for _, s2 := range arr {
		if s == s2 {
			return true
		}
	}
	return false
}

func FlagMapToList(flags map[string]string) []string {
	flagList := make([]string, 0, len(flags)*2)
	for key, value := range flags {
		key = strings.ReplaceAll(key, "_", "-")
		if !strings.HasPrefix(key, "--") {
			key = "--" + key
		}
		if value == "False" || value == "True" {
			value = strings.ToLower(value)
		} else if value != "" && value != "None" && strings.Contains(value, " ") {
			value = strconv.Quote(value)
		}
		if value != "" && value != "None" {
			flagList = append(flagList, key+"="+value)
		} else {
			flagList = append(flagList, key)
		}
	}
	return flagList
}

func GetRcloneRemotes() ([]string, error) {
	cmd := exec.Command("rclone", "listremotes")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil
	}
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(stdout)
	var remotes []string
	for scanner.Scan() {
		remotes = append(remotes, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return remotes, err
	}
	return remotes, cmd.Wait()
}

func ReplaceUnderscores(s string) string {
	sb := strings.Builder{}
	for i, r := range s {
		if i == 0 || !(r == '_' && s[i-1] == '_') {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}
