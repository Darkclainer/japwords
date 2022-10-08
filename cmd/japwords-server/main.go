package main

import (
	"fmt"
	"os"

	"japwords/pkg/config"
)

func main() {
	flagOpts := ParseFlags()
	// TODO: how to not to force config file? If it specified in
	// command line, we will complain, if not, we should search it
	// but not panic if there is no one
	// some specified parametr for config?
	userConfig, err := config.LoadConfig(flagOpts.ConfigPath, flagOpts.ConfigPathSet)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error! Failed to read config: %s\n", err)
		os.Exit(2)
	}
	app, err := NewApp(ConvertConfig(userConfig))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error! Failed to create application: %s\n", err)
		os.Exit(3)
	}
	app.Run()
}
