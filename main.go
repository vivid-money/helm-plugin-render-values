package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

var (
	debugMode = flag.Bool("debug", false, "run in debugMode")
	errLog    log.Logger
)

const (
	argPrefix = "render://"
)

type Values map[string]interface{}

type ImportValues struct {
	ImportValues []string `yaml:"importValuesFrom"`
}

func main() {
	flag.Parse()
	errLog.SetOutput(os.Stderr)

	chartArgs := os.Args
	filename := strings.TrimPrefix(chartArgs[len(chartArgs)-1], argPrefix)
	values := make(map[string]interface{})
	valuesFiles := GetValuesFiles(filename)

	values["Values"] = ReadValues(valuesFiles, filepath.Dir(filename))
	RenderTemplate(filename, values)
}

// Read file from disk
func readFile(file string) []byte {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		errLog.Fatalf("Error: \"%v\"}", err)
	}
	return data
}

// Get file list from importValuesFrom
func GetValuesFiles(file string) ImportValues {
	yamlFiles := readFile(file)
	data := ImportValues{}
	err := yaml.Unmarshal(yamlFiles, &data)
	if err != nil {
		errLog.Fatalf("Error: Can't parse file: \"%s\"; stack:\"%v\"", file, err)
	}
	if len(data.ImportValues) < 1 {
		println("There is no import values.")
	}
	return data
}

// Read values from files
func ReadValues(valuesFiles ImportValues, dir string) (vals Values) {
	vals = make(map[string]interface{})
	for _, file := range valuesFiles.ImportValues {
		yamlFiles := readFile(filepath.Join(dir, file))
		data := make(map[string]interface{})
		err := yaml.Unmarshal([]byte(yamlFiles), &data)
		if err != nil {
			errLog.Fatalf("Error: Can't parse file: \"%s\"; stack:\"%v\"", file, err)
		}
		mergeKeys(vals, data)
	}
	if len(vals) == 0 {
		vals = Values{}
	}

	return vals
}

// Recursively merge right Values into left one
func mergeKeys(left, right Values) Values {
	for key, rightVal := range right {
		if leftVal, present := left[key]; present {
			left[key] = mergeKeys(leftVal.(Values), rightVal.(Values))
		} else {
			left[key] = rightVal
		}
	}
	return left
}

// Render template to stdout
func RenderTemplate(templatefile string, data Values) {
	tpl, err := template.New(filepath.Base(templatefile)).Funcs(funcMap()).ParseFiles(templatefile)

	if err != nil {
		errLog.Fatalf("Error: Can't parse file: \"%s\"; stack:\"%v\"}", templatefile, err)
	}

	var buf strings.Builder
	err = tpl.ExecuteTemplate(&buf, filepath.Base(templatefile), data)
	if err != nil {
		errLog.Fatalf("Error: Can't render template: \"%s\"; stack:\"%v\"}", templatefile, err)
	}
	rendered := strings.ReplaceAll(buf.String(), "<no value>", "")
	if *debugMode {
		println("Values:\n---\n", rendered, "---")
	}

	fmt.Println(rendered)
}
