package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf"
	koanfyaml "github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

type SaveFailedError struct {
	Reason error
}

func (e *SaveFailedError) Error() string {
	return e.Reason.Error()
}

func (e *SaveFailedError) Unwrap() error {
	return e.Reason
}

func DefaultConfigPath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir = "./"
	}
	return filepath.Join(dir, "japwords", "config.yaml")
}

// EnsureConfigFile checks that file on path exists and if not, writes default config
func EnsureConfigFile(path string) error {
	// TODO write config
	stat, err := os.Stat(path)
	if err == nil {
		if !stat.IsDir() {
			// at least it's not a dir
			return nil
		}
		return fmt.Errorf("%s is a directory", path)
	}
	if errors.Is(err, os.ErrNotExist) {
		// we should make directories in case there are not any
		if mkdirErr := os.MkdirAll(filepath.Dir(path), 0o755); mkdirErr != nil {
			return mkdirErr
		}
		// we can write new config file
		return SaveConfig(path, DefaultUserConfig())
	}
	return err
}

func SaveConfig(path string, uc *UserConfig) error {
	buffer, err := yaml.Marshal(uc)
	if err != nil {
		return &SaveFailedError{
			Reason: err,
		}
	}
	err = os.WriteFile(path, buffer, 0o644)
	if err != nil {
		return &SaveFailedError{
			Reason: err,
		}
	}
	return nil
}

func LoadConfig(path string) (*UserConfig, error) {
	k := koanf.New(".")
	if err := k.Load(file.Provider(path), koanfyaml.Parser()); err != nil {
		return nil, fmt.Errorf("could not load file config: %s", err)
	}
	err := k.Load(env.Provider("JAPWORDS", ".", mangleEnvNames), nil)
	if err != nil {
		return nil, fmt.Errorf("could not load config from env vairables: %s", err)
	}
	var userConfig UserConfig
	unmarshalConf := koanf.UnmarshalConf{
		DecoderConfig: &mapstructure.DecoderConfig{
			ErrorUnused: true,
			Result:      &userConfig,
		},
	}
	if err := k.UnmarshalWithConf("", &userConfig, unmarshalConf); err != nil {
		return nil, fmt.Errorf("reading config failed: %s", err)
	}

	return &userConfig, nil
}

func mangleEnvNames(s string) string {
	s = strings.TrimPrefix(s, "JAPWORDS_")
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "_", ".")
	return s
}
