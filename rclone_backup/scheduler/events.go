package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

const (
	EventJobSuccessful = "rclone_backup.job_successful"
	EventJobFailed     = "rclone_backup.job_failed"
)

type EventData struct {
	Name        string  `json:"name"`
	Command     string  `json:"command"`
	Source      string  `json:"source"`
	Destination string  `json:"destination,omitempty"`
	Error       string  `json:"error,omitempty"`
	Duration    string  `json:"duration"`
	Seconds     float64 `json:"seconds"`
}

func FireEvent(type_ string, data EventData) {
	if config.NoEvents {
		return
	}

	body, err := json.Marshal(data)
	if err != nil {
		Errorln("failed to marshal event data:", err)
		return
	}

	req, err := http.NewRequest("POST", "http://supervisor/core/api/events/"+type_, bytes.NewReader(body))
	if err != nil {
		Errorln("failed to create event request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("SUPERVISOR_TOKEN"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		Errorln("failed to fire event:", err)
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		Errorln("bad status code when firing event:", resp.StatusCode)
		return
	}
}

func FireJobEvent(type_ string, job JobConfig, source string, destination string, start time.Time, msg string) {
	data := EventData{
		Name:        job.Name,
		Command:     job.Command,
		Source:      source,
		Destination: destination,
		Error:       msg,
		Duration:    FormatDuration(time.Since(start)),
		Seconds:     time.Since(start).Seconds(),
	}
	FireEvent(type_, data)
}
