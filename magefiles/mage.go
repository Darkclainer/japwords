package main

import (
	"context"
	"fmt"

	"github.com/magefile/mage/mg"
)

// Fullcheck run fullcheck for go and then for ui
func Fullcheck(ctx context.Context) error {
	fmt.Printf("Run go fullcheck\n")
	mg.SerialCtxDeps(ctx, Go.Fullcheck)
	fmt.Printf("Run ui fullcheck\n")
	mg.SerialCtxDeps(ctx, UI.Fullcheck)
	fmt.Printf("Done fullcheck!\n")
	return nil
}
