package main

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/gosimple/slug"
	"github.com/jcwillox/emerald"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	ConfigPath        = "/data/options.json"
	BackupPath        = "/backup"
	DefaultConfigPath = "/root/.config/rclone/rclone.conf"
)

var (
	AllowedSources  = []string{"/backup", "/config", "/share", "/ssl", "/media"}
	AllowedCommands = []string{"sync", "copy", "move"}
	Arrow           = emerald.Color("->", "black")
)

var (
	config   = &Config{}
	boldCyan = emerald.ColorFunc("cyan+b")
)

type Config struct {
	Jobs         []JobConfig
	ConfigPath   string `yaml:"config_path"`
	RunOnce      bool   `yaml:"run_once"`
	DryRun       bool   `yaml:"dry_run"`
	NoRename     bool   `yaml:"no_rename"`
	NoUnrename   bool   `yaml:"no_unrename"`
	LogLevel     string `yaml:"log_level"`
	RcloneConfig string `yaml:"rclone_config"`
}

type JobConfig struct {
	Name        string
	Schedule    string
	Command     string
	Sources     []string
	Destination string
	Include     []string
	Exclude     []string
	Flags       []string
}

type BackupConfig struct {
	Name string
	Slug string
}

func main() {
	// configure slug format
	slug.Lowercase = false
	slug.CustomSub = map[string]string{
		" ": "_",
		"(": "",
		")": "",
		"[": "",
		"]": "",
		":": "",
		",": "",
	}

	// load addon configuration
	var err error
	config, err = LoadConfig()
	if err != nil {
		Fatalln("failed to read or parse config", err)
	}

	// write rclone config from addon config
	if config.RcloneConfig != "" {
		err := os.WriteFile(DefaultConfigPath, []byte(config.RcloneConfig), 0666)
		if err != nil {
			Fatalln("failed to write config")
		}
		config.ConfigPath = DefaultConfigPath
	}

	// check rclone config exists
	if stat, _ := os.Stat(config.ConfigPath); stat == nil {
		Fatalln(
			"rclone config does not exist!" +
				"\nIf this is your first time starting this add-on ensure to" +
				"\ncreate a valid rclone configuration at \"" + config.ConfigPath + "\"",
		)
	} else {
		Infoln("rclone config found")
	}

	Infoln("checking job configs...")
	for _, job := range config.Jobs {
		err := CheckJob(job)
		if err != nil {
			Fatalln(err)
		}
	}

	if config.RunOnce {
		for _, job := range config.Jobs {
			if job.Schedule == "" {
				CreateJob(job)()
			}
		}
	} else {
		scheduler := gocron.NewScheduler(time.Local)

		// only run 1 job at a time to prevent issues with file locks
		scheduler.SetMaxConcurrentJobs(1, gocron.WaitMode)

		Infoln("scheduled jobs:")

		for _, job := range config.Jobs {
			if job.Schedule == "" {
				// schedule to run immediately
				_, err = scheduler.Every(1).Second().LimitRunsTo(1).Do(CreateJob(job))
			} else {
				_, err = scheduler.Cron(job.Schedule).Do(CreateJob(job))
			}

			if err != nil {
				Fatalln("failed to schedule job", "'"+job.Name+"'", err)
			}

			PrintJob(job)
		}

		scheduler.StartBlocking()
	}
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile(ConfigPath)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func CheckJob(job JobConfig) error {
	// check allowed command
	if !ArrayContains(AllowedCommands, job.Command) {
		return fmt.Errorf("command '%s' is not allowed; must be one of %s", job.Command, AllowedCommands)
	}
	for _, source := range job.Sources {
		// check allowed source
		if !ArrayHasPrefix(AllowedSources, source) {
			return fmt.Errorf("source '%s' is not allowed; must be one of %s", source, AllowedSources)
		}
		// check source exists
		if stat, err := os.Stat(source); stat == nil {
			return fmt.Errorf("source '%s' does not exist; %v", source, err)
		}
	}
	return nil
}

// CreateJob create run job closure with the job config
func CreateJob(job JobConfig) func() {
	return func() {
		for _, source := range job.Sources {
			// adjust destination with multiple sources
			subfolder := ""
			if len(job.Sources) > 1 {
				subfolder = source
			}

			// generate rclone command
			args := []string{job.Command, source, job.Destination + subfolder, "--verbose", "--config", config.ConfigPath}

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
			args = append(args, job.Flags...)

			start := time.Now()

			cmd := exec.Command("rclone", args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin

			if job.Name != "" {
				Infoln("running", "\""+emerald.Cyan+job.Name+emerald.Green+"\";", emerald.HighlightPathStat(source), Arrow, job.Destination)
			} else {
				Infoln("running job", emerald.HighlightPathStat(source), Arrow, job.Destination)
			}

			Debugln("rclone", args)

			var undoRename func()
			if strings.HasPrefix(source, BackupPath) && !config.NoRename {
				var err error
				undoRename, err = RenameBackups()
				if err != nil {
					Errorln("failed to rename backups, aborting upload", err)
					return
				}
			}

			err := cmd.Run()
			if err != nil {
				Errorln("failed to run rclone command", err)
				return
			}

			if undoRename != nil && !config.NoUnrename {
				undoRename()
			}

			Infoln("finished in", boldCyan(FormatDuration(time.Since(start))))
		}
	}
}
