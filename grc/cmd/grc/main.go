package main

import (
	"os"

	"github.com/carlosrabelo/grc/grc/internal/app"
)

func main() {
	os.Exit(app.Run(os.Args[1:]))
}
