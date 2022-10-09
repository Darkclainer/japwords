package main

import (
	"os"

	"github.com/magefile/mage/sh"
)

func RunVDir(dir string, cmd string, args ...string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := os.Chdir(dir); err != nil {
		return err
	}
	defer os.Chdir(wd)
	return sh.RunV(cmd, args...)
}
