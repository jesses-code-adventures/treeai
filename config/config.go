package config

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
)

type Config struct {
	Bin       string
	Commands  []string
	Data      string
	Debug     bool
	Silent    bool
	Gitignore bool
	Window    bool
}

func New() *Config {
	return &Config{
		Bin:       "opencode",
		Commands:  []string{},
		Data:      os.ExpandEnv("$HOME/.local/share/treeai"),
		Debug:     false,
		Silent:    false,
		Gitignore: false,
		Window:    false,
	}
}

func (c *Config) WorktreePath(worktreeName string) string {
	return filepath.Join(c.Data, worktreeName)
}

func (c *Config) ToSlogAttrs() []any {
	data, _ := json.Marshal(c)

	var m map[string]any
	err := json.Unmarshal(data, &m)
	if err != nil {
		panic(err)
	}

	attrs := make([]any, 0, len(m)*2)
	for k, v := range m {
		attrs = append(attrs, k, v)
	}

	return attrs
}

func (c *Config) ApplyFlags(bin string, silent bool, data string, windowCommands []string, useGitignore, debug, window bool) {
	// only override if flag was explicitly set (you'll need to track this in cobra)
	if bin != "opencode" {
		c.Bin = bin
	}
	if silent {
		c.Silent = silent
	}
	if useGitignore {
		c.Gitignore = useGitignore
	}
	if len(windowCommands) > 0 {
		c.Commands = windowCommands
	}
	if debug {
		c.Debug = debug
	}
	if window {
		c.Window = window
	}
	if data != "" {
		c.Data = data
	}
}

func Load() (*Config, error) {
	cfg := New()

	configPath := getConfigPath()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return cfg, nil // return defaults if no config file
	}

	if _, err := toml.DecodeFile(configPath, cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func getConfigPath() string {
	if configDir := os.Getenv("XDG_CONFIG_HOME"); configDir != "" {
		return filepath.Join(configDir, "treeai", "config.toml")
	}

	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "treeai", "config.toml")
}
