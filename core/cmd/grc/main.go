package main

import (
	"context"
	"fmt"
	"os"

	"github.com/carlosrabelo/grc/core/internal/app"
)

var (
	version   = "dev"
	buildTime = "unknown"
)

func main() {
	ctx := context.Background()
	if err := app.Run(ctx, version, buildTime, os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "grc: %v\n", err)
		os.Exit(1)
	}
}
