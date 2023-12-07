package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Run mg.Namespace

// UI run ui development server
func (Run) UI() error {
	return UIRunV(
		"npm",
		"run",
		"dev",
	)
}

// UI run backend development server
func (Run) Server() error {
	return sh.RunV(
		"go",
		"run",
		"./cmd/japwords-server",
	)
}
