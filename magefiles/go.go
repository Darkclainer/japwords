package main

import (
	"context"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Go mg.Namespace

func (Go) Lint(ctx context.Context) error {
	mg.CtxDeps(ctx, mg.F(Tools.Install, "golangci-lint"))
	golangci := getToolPath("golangci-lint")
	return sh.RunV(golangci, "run", "--sort-results")
}

func (Go) Generate(ctx context.Context) error {
	mg.CtxDeps(ctx, mg.F(Tools.Install, "mockery"))
	mockery, err := filepath.Abs(getToolPath("mockery"))
	if err != nil {
		return err
	}
	enumer, err := filepath.Abs(getToolPath("enumer"))
	if err != nil {
		return err
	}
	return sh.RunWithV(
		map[string]string{
			"MOCKERY_TOOL": mockery,
			"ENUMER_TOOL":  enumer,
		},
		"go", "generate", "./...",
	)
}
