package main

import (
	"context"
	"fmt"
	"os"

	"github.com/carlosrabelo/grc/core/internal/app"
)

// Build information injected via ldflags
var (
	version   = "dev"
	buildTime = "unknown"
)

func main() {
	ctx := context.Background()
	
	if err := runApplication(ctx); err != nil {
		handleError(err)
	}
}

// runApplication runs the application with the provided context
func runApplication(ctx context.Context) error {
	return app.Run(ctx, version, buildTime, os.Args[1:], os.Stdout, os.Stderr)
}

// handleError handles errors consistently
func handleError(err error) {
	fmt.Fprintf(os.Stderr, "grc: %v\n", err)
	os.Exit(1)
}
