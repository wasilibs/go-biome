package main

import (
	"os"

	"github.com/wasilibs/go-biome/v2/internal/runner"
)

func main() {
	os.Exit(runner.Run("biome", os.Args[1:], os.Stdin, os.Stdout, os.Stderr, "."))
}
