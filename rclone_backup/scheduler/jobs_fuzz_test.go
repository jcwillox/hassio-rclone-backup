package main

import (
	"os"
	"regexp"
	"slices"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// FuzzFlagMapToList each flag becomes exactly one `--name[=value]` argument, and the value reaches
// rclone as written: there is no shell to remove quotes, so rclone would read them as part of it
func FuzzFlagMapToList(f *testing.F) {
	f.Add("drive_use_trash", "False")
	f.Add("min-age", "14d")
	f.Add("--fast-list", "")
	f.Add("exclude", "None")
	f.Add("user-agent", "my agent")
	f.Add("include", "a b.txt")
	f.Add("header", `X-Note: "quoted"`)
	f.Fuzz(func(t *testing.T, key string, value string) {
		args := FlagMapToList(Flags{key: value})
		if len(args) != 1 {
			t.Fatalf("flag %q=%q gave %d arguments: %q", key, value, len(args), args)
		}
		arg := args[0]
		if !strings.HasPrefix(arg, "--") {
			t.Fatalf("flag %q gave %q, which does not start with --", key, arg)
		}
		// underscores stand for dashes, also in a leading --
		name := strings.TrimPrefix(strings.ReplaceAll(key, "_", "-"), "--")
		if value == "" || value == "None" {
			if arg != "--"+name {
				t.Fatalf("flag %q with no value gave %q", key, arg)
			}
			return
		}
		want := value
		if value == "True" || value == "False" {
			want = strings.ToLower(value)
		}
		if arg != "--"+name+"="+want {
			t.Fatalf("flag %q=%q gave %q, want %q", key, value, arg, "--"+name+"="+want)
		}
	})
}

// FuzzBuildArgs the command, paths and filters are always in place, and dry_run adds --dry-run
func FuzzBuildArgs(f *testing.F) {
	f.Add("sync", "/backup", "google:/Backup", "DailyBackup*", "", false)
	f.Add("copy", "/share", "", "", "*.tmp", true)
	f.Add("move", "/media", "box:", "--dry-run", "--include", false)
	f.Fuzz(func(t *testing.T, command, source, destination, include, exclude string, dryRun bool) {
		saved := config
		t.Cleanup(func() { config = saved })
		config = &Config{DryRun: dryRun}

		job := JobConfig{Command: command}
		if include != "" {
			job.Include = []string{include}
		}
		if exclude != "" {
			job.Exclude = []string{exclude}
		}
		args := BuildArgs(job, source, destination)

		want := []string{command, source}
		if destination != "" {
			want = append(want, destination)
		}
		want = append(want, "--verbose")
		if include != "" {
			want = append(want, "--include", include)
		}
		if exclude != "" {
			want = append(want, "--exclude", exclude)
		}
		if dryRun {
			want = append(want, "--dry-run")
		}
		if !slices.Equal(args, want) {
			t.Fatalf("got %q, want %q", args, want)
		}
	})
}

type addonConfig struct {
	Options struct {
		Jobs []struct {
			Command string
		}
	}
	Schema struct {
		Jobs []struct {
			Command string
		}
	}
}

// commandSchema the allowed job commands, read from the add-on's config.yaml schema
func commandSchema(t testing.TB) []string {
	data, err := os.ReadFile("../config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	var addon addonConfig
	if err := yaml.Unmarshal(data, &addon); err != nil {
		t.Fatal(err)
	}
	// Supervisor's list(val1|val2) type: an exact match on one of the values
	// (supervisor/apps/options.py, RE_SCHEMA_ELEMENT and vol.In)
	match := regexp.MustCompile(`^list\((.+)\)\??$`).FindStringSubmatch(addon.Schema.Jobs[0].Command)
	if match == nil {
		t.Fatalf("job command schema %q is not a list", addon.Schema.Jobs[0].Command)
	}
	return strings.Split(match[1], "|")
}

// FuzzCommandSchema only copy, sync, move and delete are accepted as a job command
func FuzzCommandSchema(f *testing.F) {
	allowed := commandSchema(f)
	for _, seed := range []string{"copy", "sync", "move", "ls", "delete", "Sync", " sync", "copy|sync", ""} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, command string) {
		accepted := slices.Contains(allowed, command)
		want := command == "copy" || command == "sync" || command == "move" || command == "delete"
		if accepted != want {
			t.Fatalf("command %q: accepted %v, want %v", command, accepted, want)
		}
	})
}

// TestDocumentedCommands every job command in the default options and the docs is accepted
func TestDocumentedCommands(t *testing.T) {
	allowed := commandSchema(t)

	data, err := os.ReadFile("../config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	var addon addonConfig
	if err := yaml.Unmarshal(data, &addon); err != nil {
		t.Fatal(err)
	}
	for _, job := range addon.Options.Jobs {
		if !slices.Contains(allowed, job.Command) {
			t.Errorf("config.yaml: default job command %q is not allowed by the schema", job.Command)
		}
	}

	command := regexp.MustCompile(`(?m)^\s*-?\s*command:\s*(\S+)`)
	for _, path := range []string{"../DOCS.md", "../../README.md"} {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		found := command.FindAllStringSubmatch(string(data), -1)
		if len(found) == 0 {
			t.Fatalf("%s: no job command found", path)
		}
		for _, match := range found {
			if !slices.Contains(allowed, match[1]) {
				t.Errorf("%s: job command %q is not allowed by the schema", path, match[1])
			}
		}
	}
}
