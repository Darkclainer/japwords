package main

import (
	"fmt"
	"os"

	"github.com/Darkclainer/japwords/cmd/japwords-server/fxapp"
	"github.com/Darkclainer/japwords/pkg/config"
)

func main() {
	flagOpts := ParseFlags()
	configMgr, err := prepareConfig(flagOpts.ConfigPath, flagOpts.ConfigPathSet)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error! Failed to read config: %s\n", err)
		os.Exit(2)
	}
	app, err := fxapp.NewApp(configMgr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error! Failed to create application: %s\n", err)
		os.Exit(3)
	}
	app.Run()
}

func prepareConfig(path string, provided bool) (*config.Manager, error) {
	if !provided {
		path = config.DefaultConfigPath()
	}
	err := config.EnsureConfigFile(path)
	if err != nil {
		return nil, err
	}
	return config.New(path)
}
