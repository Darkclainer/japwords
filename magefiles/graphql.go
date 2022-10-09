package main

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Graphql mg.Namespace

// Generate regenerates server graphql code
func (Graphql) GenerateServer(ctx context.Context) error {
	mg.CtxDeps(ctx, Tools.Gqlgen)

	fmt.Printf("Generating server graphql files\n")
	defer fmt.Printf("Done server graphql files generation\n")
	return sh.RunV(
		getToolPath("gqlgen"),
		"--config",
		filepath.Join("gqlgen.yaml"),
		"generate",
	)
}

// Generate regenerates ui graphql code
func (Graphql) GenerateUI(ctx context.Context) error {
	return RunVDir(
		"ui",
		"npm",
		"run",
		"generate",
	)
}

// Generate regenerate server and ui graphql code
func (Graphql) Generate(ctx context.Context) {
	mg.SerialCtxDeps(
		ctx,
		Graphql.GenerateServer,
		Graphql.GenerateUI,
	)
}
