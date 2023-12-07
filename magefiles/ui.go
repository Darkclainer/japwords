package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/magefile/mage/mg"
)

type UI mg.Namespace

const uiDir = "ui"

func UIRunV(cmd string, args ...string) error {
	nvmDir := os.Getenv("NVM_DIR")
	if nvmDir != "" {
		nvmExec := filepath.Join(nvmDir, "nvm-exec")
		return RunVDir(
			uiDir,
			nvmExec,
			append([]string{cmd}, args...)...,
		)
	}
	fmt.Println("NVM configuration not found, trying to run command as is!")
	return RunVDir(
		uiDir,
		cmd,
		args...,
	)
}

func (UI) Lint(ctx context.Context) error {
	return UIRunV(
		"npm",
		"run",
		"lint",
	)
}

func (UI) Typecheck(ctx context.Context) error {
	return UIRunV(
		"npm",
		"run",
		"type-check",
	)
}

func (UI) Test(ctx context.Context) error {
	return UIRunV(
		"npm",
		"run",
		"test",
	)
}

// Fullcheck generates go code, then generates graphql, then runs tests and lint
func (UI) Fullcheck(ctx context.Context) error {
	fmt.Printf("Run graphql generate ui\n")
	mg.SerialCtxDeps(ctx, Graphql.GenerateUI)
	fmt.Printf("Run type check\n")
	mg.SerialCtxDeps(ctx, UI.Typecheck)
	fmt.Printf("Run tests\n")
	mg.SerialCtxDeps(ctx, UI.Test)
	fmt.Printf("Run lint\n")
	mg.SerialCtxDeps(ctx, UI.Lint)
	fmt.Printf("Done!\n")
	return nil
}
