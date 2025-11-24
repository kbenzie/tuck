package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"tuck/internal/path"

	"go.yaml.in/yaml/v4"
)

var (
	ConfigFile = filepath.Join(path.ConfigDir, "tuck.yaml")
)

// The filters below are used to select the release assets based on properties
// of the localhost; required properties such as operating system and CPU
// architecture, these must all match for an asset to be selected; optional
// properties such as the linked C standard library, these will be used in the
// event there are multiple candiate releases assets to choose from.
type ConfigFilters struct {
	Required []string `yaml:"required"`
	Optional []string `yaml:"optional"`
}

type Config struct {
	Filters ConfigFilters `yaml:"filters"`
}

func detectArchFilter() string {
	switch runtime.GOARCH {
	case "amd64":
		return "(amd64|x86-64|x86_64)"
	case "arm64":
		return "(arm64|aarch64)"
	default:
		// TODO: Handle other architectures
		log.Fatalln("unimplemented arch:", runtime.GOARCH)
		return ""
	}
}

func linuxDefaultFilters() ConfigFilters {
	filters := ConfigFilters{}
	filters.Required = append(filters.Required,
		"linux",
		"(.tar.(gz|bz2|xz)|.zip)$",
	)
	filters.Optional = append(filters.Optional,
		detectArchFilter(),
		"musl",
	)
	return filters
}

func darwinDefaultFilters() ConfigFilters {
	filters := ConfigFilters{}
	filters.Required = append(filters.Required,
		"(mac|macos|darwin)",
		"(.tar.(gz|bz2|xz)|.zip)$",
	)
	filters.Optional = append(filters.Optional,
		detectArchFilter(),
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
		case "darwin":
			config.Filters = darwinDefaultFilters()
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
