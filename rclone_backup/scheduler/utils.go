package main

import (
	"github.com/jcwillox/emerald"
	"strconv"
	"strings"
	"time"
)

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

func ArrayHasPrefix(arr []string, s string) bool {
	for _, s2 := range arr {
		if strings.HasPrefix(s, s2) {
			return true
		}
	}
	return false
}

func PrintJob(job JobConfig) {
	if job.Schedule != "" {
		emerald.Print(job.Schedule, " ")
	}
	if job.Name != "" {
		emerald.Print(emerald.Red, "\"", job.Name, "\" ", emerald.Reset)
	}
	emerald.Print("[")
	for i, source := range job.Sources {
		emerald.Print(emerald.HighlightPathStat(source))
		if i < len(job.Sources)-1 {
			emerald.Print(",")
		}
	}
	emerald.Print("] ", job.Destination, "\n")
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
