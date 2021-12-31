// Package config contains code for config parsing
package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/microhod/repo/internal/path"
)

const (
	configFolder = "~/.config/repo"
	configFile   = "config.json"
)

var (
	ConfigFolder = path.Clean(configFolder)
	ConfigFile   = path.Clean(fmt.Sprintf("%s/%s", configFolder, configFile))

	defaultConfig = Config{
		Remote: RemoteConfig{
			Default: DefaultRemoteConfig{
				Prefix: "ssh://git@github.com",
			},
		},
		Local: LocalConfig{
			Root: "~/src",
		},
	}
)

type Config struct {
	Remote RemoteConfig `json:"remote"`
	Local  LocalConfig  `json:"local"`
}

type RemoteConfig struct {
	Default DefaultRemoteConfig `json:"default"`
}

type DefaultRemoteConfig struct {
	Prefix string `json:"prefix"`
}

type LocalConfig struct {
	Root string `json:"root"`
}

func Parse() (Config, error) {
	if err := writeDefaultConfig(); err != nil {
		return Config{}, err
	}
	file, err := os.ReadFile(ConfigFile)
	if err != nil {
		return Config{}, err
	}
	cfg := &Config{}
	if err := json.Unmarshal(file, cfg); err != nil {
		return Config{}, err
	}
	return *cfg, nil
}

func writeDefaultConfig() error {
	if _, err := os.Stat(ConfigFile); !os.IsNotExist(err) {
		return nil
	}
	if err := os.MkdirAll(ConfigFolder, os.ModePerm); err != nil {
		return err
	}
	file, err := os.Create(ConfigFile)
	if err != nil {
		return err
	}

	cfg, err := json.MarshalIndent(defaultConfig, "", "    ")
	if err != nil {
		return err
	}
	_, err = file.Write(cfg)
	return err
}
