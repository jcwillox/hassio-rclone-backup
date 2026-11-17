package main

import (
	"slices"
	"strings"
	"testing"
)

// TestCheckJob each command gets the paths it needs, and a delete never runs without a filter
func TestCheckJob(t *testing.T) {
	saved := remotes
	t.Cleanup(func() { remotes = saved })
	remotes = []string{"hassbackup:", "google:"}

	tests := []struct {
		name string
		job  JobConfig
		err  string
	}{
		{"sync with destination", JobConfig{Command: "sync", Sources: []string{"/"}, Destinations: []string{"google:/Backup"}}, ""},
		{"copy without destination", JobConfig{Command: "copy", Sources: []string{"/"}}, "needs a destination"},
		{"no source", JobConfig{Command: "delete"}, "at least 1 source"},
		{"delete with min-age", JobConfig{Command: "delete", Sources: []string{"hassbackup:"}, Flags: Flags{"min-age": "14d"}}, ""},
		{"delete with max_age", JobConfig{Command: "delete", Sources: []string{"hassbackup:"}, Flags: Flags{"max_age": "1y"}}, ""},
		{"delete with --min-age key", JobConfig{Command: "delete", Sources: []string{"hassbackup:"}, Flags: Flags{"--min-age": "14d"}}, ""},
		{"delete with include", JobConfig{Command: "delete", Sources: []string{"hassbackup:"}, Include: []string{"Old*"}}, ""},
		{"delete with exclude", JobConfig{Command: "delete", Sources: []string{"hassbackup:"}, Exclude: []string{"Keep*"}}, ""},
		{"delete with extra flag", JobConfig{Command: "delete", Sources: []string{"hassbackup:"}, ExtraFlags: []string{"--min-age=14d"}}, ""},
		{"delete with split extra flag", JobConfig{Command: "delete", Sources: []string{"hassbackup:", "google:"}, ExtraFlags: []string{"--min-age", "14d"}}, ""},
		{"delete with no filter", JobConfig{Command: "delete", Sources: []string{"hassbackup:"}}, "needs a filter"},
		{"delete with other flags only", JobConfig{Command: "delete", Sources: []string{"hassbackup:"}, Flags: Flags{"drive-use-trash": "False"}}, "needs a filter"},
		{"delete with destination", JobConfig{Command: "delete", Sources: []string{"hassbackup:"}, Destinations: []string{"google:"}, Flags: Flags{"min-age": "14d"}}, "takes no destination"},
		{"delete of a local path", JobConfig{Command: "delete", Sources: []string{"/"}, Flags: Flags{"min-age": "14d"}}, "must be a remote"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := CheckJob(test.job)
			if test.err == "" {
				if err != nil {
					t.Fatalf("got error %q, want none", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), test.err) {
				t.Fatalf("got error %v, want one containing %q", err, test.err)
			}
		})
	}
}

// TestJobRuns a job with no destination never uses a source as its destination
func TestJobRuns(t *testing.T) {
	tests := []struct {
		name string
		job  JobConfig
		want [][2]string
	}{
		{"delete one source", JobConfig{Sources: []string{"hassbackup:"}}, [][2]string{{"hassbackup:", ""}}},
		{"delete two sources", JobConfig{Sources: []string{"hassbackup:", "google:/Backup"}}, [][2]string{{"hassbackup:", ""}, {"google:/Backup", ""}}},
		{"one to one", JobConfig{Sources: []string{"/backup"}, Destinations: []string{"google:"}}, [][2]string{{"/backup", "google:"}}},
		{"two sources", JobConfig{Sources: []string{"/share", "/ssl"}, Destinations: []string{"google:"}}, [][2]string{{"/share", "google:/share"}, {"/ssl", "google:/ssl"}}},
		{"two destinations", JobConfig{Sources: []string{"/backup"}, Destinations: []string{"google:", "box:"}}, [][2]string{{"/backup", "google:"}, {"/backup", "box:"}}},
		{"two by two", JobConfig{Sources: []string{"/share", "/ssl"}, Destinations: []string{"google:", "box:"}}, [][2]string{{"/share", "google:/share"}, {"/ssl", "google:/ssl"}, {"/share", "box:/share"}, {"/ssl", "box:/ssl"}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := JobRuns(test.job); !slices.Equal(got, test.want) {
				t.Fatalf("got %q, want %q", got, test.want)
			}
		})
	}
}
