package wasm

import _ "embed"

//go:embed biome.wasm
var Biome []byte

//go:embed memory.wasm
var Memory []byte
