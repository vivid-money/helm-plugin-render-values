package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var (
	debugMode = flag.Bool("debug", false, "run in debugMode")
	errLog    log.Logger
)

const (
	argPrefix = "render://"
)

func main() {
	flag.Parse()
	args := flag.Args()
	errLog.SetOutput(os.Stderr)
	if len(args) == 0 {
		errLog.Fatalf(`
		ERROR: Required arguments are missing.
		Usage: %s %s<values-filename.yaml>`, filepath.Base(os.Args[0]), argPrefix)
	}
	argFilename := args[len(args)-1]
	valuesRender := new(ValuesRenderer)
	valuesRender.Run(argFilename, argPrefix)
}

// Read the file.
func readFile(file string) []byte {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		errLog.Fatalf("ERROR: \"%v\"}", err)
	}
	return data
}

// Recursively merge right Values into left one
func mergeKeys(left, right Values) Values {
	for key, rightVal := range right {
		if leftVal, present := left[key]; present {
			if _, ok := leftVal.(Values); ok {
				left[key] = mergeKeys(leftVal.(Values), rightVal.(Values))
			} else {
				left[key] = rightVal
			}
		} else {
			left[key] = rightVal
		}
	}
	return left
}

// If glob is used. Then map folder tree as yaml structure.
func DirrectoryMapping(dir []string, val Values) Values {
	data := make(Values)
	if len(dir) > 1 {
		data[dir[0]] = DirrectoryMapping(dir[1:], val)
	} else {
		data[dir[0]] = val
	}
	return data
}

// Make file path absolut depend on main file.
func absolutePath(basefile, depenfile string) string {
	dir := filepath.Dir(basefile)
	return filepath.Join(dir, depenfile)
}
