package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"tuck/internal/path"

	"go.yaml.in/yaml/v4"
)

var (
	ConfigFile = filepath.Join(path.ConfigDir, "tuck.yaml")
)

type ConfigFilters struct {
	Required []string `yaml:"required"`
	Optional []string `yaml:"optional"`
}

type Config struct {
	Filters ConfigFilters `yaml:"filters"`
}

func linuxDefaultFilters() ConfigFilters {
	filters := ConfigFilters{}
	filters.Required = append(filters.Required,
		"linux",
		"(amd64|x86-64|x86_64)", // TODO: detect this
		".tar.(gz|bz2|xz)",
	)
	filters.Optional = append(filters.Optional,
		"musl",
	)
	return filters
}

func Load() (Config, error) {
	config := Config{}
	if path.Exists(ConfigFile) {
		data, err := os.ReadFile(ConfigFile)
		if err != nil {
			return config, err
		}
		yaml.Unmarshal(data, &config)
		// TODO: how to handle updating default filters?
	} else {
		switch runtime.GOOS {
		case "linux":
			config.Filters = linuxDefaultFilters()
		default:
			return config, fmt.Errorf("unimplemented OS: %s", runtime.GOOS)
		}
	}
	return config, nil
}

func Store(config Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigFile, data, 0644)
}
