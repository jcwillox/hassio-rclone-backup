package main

import (
	"bufio"
	"os/exec"
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
	cmd := exec.Command("rclone", "listremotes", "--config", config.ConfigPath)
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
