package main

import (
	"flag"
	"fmt"
	"os"
)

type FlagOpts struct {
	ConfigPath string
	// ConfigPathSet indicates that path for config was provided by user
	ConfigPathSet bool
}

func ParseFlags() *FlagOpts {
	fset := newFlagSet()

	// sloppy, but ok
	var flagOpts FlagOpts
	fset.StringVar(&flagOpts.ConfigPath, "c", "config.yaml", "path to config")
	fset.Parse(os.Args[1:])

	fset.Visit(func(f *flag.Flag) {
		if f.Name == "c" {
			flagOpts.ConfigPathSet = true
		}
	})

	return &flagOpts
}

func newFlagSet() *flag.FlagSet {
	cliName := "japwords"
	fset := flag.NewFlagSet(cliName, flag.ExitOnError)
	fset.SetOutput(os.Stderr)
	fset.Usage = func() { printUsage(fset, cliName) }
	return fset
}

func printUsage(f *flag.FlagSet, name string) {
	fmt.Fprintf(f.Output(), "Usage:\n  %s [flags]\n", name)

	// print flags
	fmt.Fprint(f.Output(), "\nflags:\n")
	f.PrintDefaults()
}
