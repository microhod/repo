package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/microhod/repo/internal/path"
)

var (
	configFolder = path.Clean("~/.config/repo")
	configFile   = path.Clean(fmt.Sprintf("%s/config.json", configFolder))

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
		return Config{}, fmt.Errorf("writing default cfg: %w", err)
	}

	f, err := os.Open(configFile)
	if err != nil {
		return Config{}, fmt.Errorf("open cfg: %w", err)
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("unmarshal: %w", err)
	}
	return cfg, nil
}

func writeDefaultConfig() error {
	_, err := os.Stat(configFile)
	if err == nil {
		// exit early if a config file already exists
		return nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("checking cfg file exists: %w", err)
	}

	if err := os.MkdirAll(configFolder, os.ModePerm); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}
	cfg, err := json.MarshalIndent(defaultConfig, "", "    ")
	if err != nil {
		return fmt.Errorf("marshal cfg: %w", err)
	}
	if err := os.WriteFile(configFile, cfg, 0644); err != nil {
		return fmt.Errorf("writing cfg: %w", err)
	}
	return nil
}
