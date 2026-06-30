package main

import (
	"github.com/goyek/x/boot"
	"github.com/wasilibs/tools/tasks"
)

func main() {
	tasks.Define(tasks.Params{
		LibraryName: "biome",
		LibraryRepo: "biomejs/biome",
		GoReleaser:  true,
	})
	boot.Main()
}
