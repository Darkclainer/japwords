package main

import (
	"context"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Graphql mg.Namespace

// Generate regenerates graphql go code
func (Graphql) Generate(ctx context.Context) error {
	mg.Deps(Tools.Gqlgen)

	return sh.RunV(
		getToolPath("gqlgen"),
		"--config",
		filepath.Join("gqlgen.yaml"),
		"generate",
	)
}
