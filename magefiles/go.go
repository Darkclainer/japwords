package main

import (
	"context"
	"fmt"
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

func (Go) Test(ctx context.Context) error {
	return sh.RunV(
		"go", "test", "./...",
	)
}

// Fullcheck generates go code, then graphql, then runs tests and lint
func (Go) Fullcheck(ctx context.Context) error {
	fmt.Printf("Run go generate\n")
	mg.SerialCtxDeps(ctx, Go.Generate)
	fmt.Printf("Run graphql generate server\n")
	mg.SerialCtxDeps(ctx, Graphql.GenerateServer)
	fmt.Printf("Run go test\n")
	mg.SerialCtxDeps(ctx, Go.Test)
	fmt.Printf("Run go lint\n")
	mg.SerialCtxDeps(ctx, Go.Lint)
	fmt.Printf("Done!\n")
	return nil
}
