package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/mitchellh/mapstructure"
)

func LoadConfig(path string, provided bool) (*UserConfig, error) {
	k := koanf.New(".")
	if path != "" {
		if err := k.Load(optionalFileProvider(path, provided), yaml.Parser()); err != nil {
			return nil, fmt.Errorf("could not load file config: %s", err)
		}
	}
	k.Load(env.Provider("JAPWORDS", ".", mangleEnvNames), nil)
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

// optionalFileProvider return file provider that do not generate error
// if provided is false and path does not exists.
func optionalFileProvider(path string, provided bool) *optionalFile {
	return &optionalFile{
		path:     filepath.Clean(path),
		provided: provided,
	}
}

type optionalFile struct {
	path string
	// provided specify should not existing path generate error
	provided bool
}

func (f *optionalFile) ReadBytes() ([]byte, error) {
	src, err := os.ReadFile(f.path)
	if err != nil {
		if !f.provided && errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	return src, nil
}

func (f *optionalFile) Read() (map[string]interface{}, error) {
	return nil, errors.New("optional file provider does not support this method")
}

func (f *optionalFile) Watch(cb func(event interface{}, err error)) error {
	return errors.New("optional file provider does not support this method")
}
