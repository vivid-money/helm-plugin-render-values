package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"main/render"
)

var (
	debugMode = flag.Bool("debug", false, "run in debugMode")
)

const (
	argPrefix = "render://"
)

func main() {
	var errLog log.Logger
	errLog.SetOutput(os.Stderr)

	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		errLog.Fatalf(`
		ERROR: Required arguments are missing.
		Usage: %s %s<values-filename.yaml>`, filepath.Base(os.Args[0]), argPrefix)
	}
	argFilename := args[len(args)-1]
	valuesRender := new(render.ValuesRenderer)
	valuesRender.Debug = *debugMode
	valuesRender.Run(argFilename, argPrefix)
}
