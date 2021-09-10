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

type FilesList struct {
	ImportValues []string `yaml:"importValuesFrom"`
	ExtendRender []string `yaml:"extendRenderWith"`
}

type ValuesRenderer struct {
	filename string
	files    FilesList
	values   Values
}

func main() {
	flag.Parse()
	errLog.SetOutput(os.Stderr)

	chartArg := os.Args[len(os.Args)-1]
	vlRender := new(ValuesRenderer)
	vlRender.filename = strings.TrimPrefix(chartArg, argPrefix)

	vlRender.GetFileList()
	vlRender.ReadValues()
	vlRender.RenderTemplate()
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
func (vr *ValuesRenderer) GetFileList() error {
	dir := filepath.Dir(vr.filename)
	rawFile := readFile(vr.filename)
	yamlFiles := strings.ReplaceAll(string(rawFile), "{{", "#{{")
	data := FilesList{}
	err := yaml.Unmarshal([]byte(yamlFiles), &data)
	if err != nil {
		errLog.Fatalf("Error GetValuesFiles: Can't parse file: \"%s\"; stack:\"%v\"", vr.filename, err)
	}
	if len(data.ImportValues) < 1 {
		println("Info: there is no import values.")
	} else {
		for i, source := range data.ImportValues {
			if source == "self" {
				data.ImportValues[i] = vr.filename
				println("Info: there is itself using for values.")
			} else {
				data.ImportValues[i] = filepath.Join(dir, source)
			}
		}
	}
	if len(data.ExtendRender) < 1 {
		println("Info: there is no extended files for render.")
	} else {
		for i, source := range data.ExtendRender {
			data.ExtendRender[i] = filepath.Join(dir, source)
		}
	}

	data.ExtendRender = append(data.ExtendRender, vr.filename)
	vr.files = data

	return err
}

// Read values from files
func (vr *ValuesRenderer) ReadValues() {

	vals := make(Values)

	for _, file := range vr.files.ImportValues {
		rawFile := readFile(file)
		yamlFiles := strings.ReplaceAll(string(rawFile), "{{", "#{{")
		data := make(Values)
		err := yaml.Unmarshal([]byte(yamlFiles), &data)
		if err != nil {
			errLog.Fatalf("Error ReadValues: Can't parse file: \"%s\"; stack:\"%v\"", file, err)
		}
		mergeKeys(vals, data)
	}
	if len(vals) == 0 {
		vals = Values{}
	}
	vr.values = make(Values)
	vr.values["Values"] = vals
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

// Render template to stdout
func (vr *ValuesRenderer) RenderTemplate() {
	valuesResult := make(Values)
	for _, file := range vr.files.ExtendRender {
		tpl, err := template.New("render").Funcs(funcMap()).ParseFiles(file)

		if err != nil {
			errLog.Fatalf("Error create render: stack:\"%v\"}", err)
		}

		var buf strings.Builder

		err = tpl.ExecuteTemplate(&buf, filepath.Base(file), vr.values)
		if err != nil {
			errLog.Fatalf("Error: Can't render template: \"%s\"; stack:\"%v\"}", file, err)
		}
		rendered := strings.ReplaceAll(buf.String(), "<no value>", "")

		data := make(Values)
		err = yaml.Unmarshal([]byte(rendered), &data)
		if err != nil {
			errLog.Fatalf("Error ReadValues: Can't parse file: \"%s\"; stack:\"%v\"", file, err)
		}
		mergeKeys(valuesResult, data)

	}
	renderedValues, err := yaml.Marshal(valuesResult)
	if err != nil {
		log.Fatalf("Error: Can't execute toYAML func:\"%v\"\n   \"%s\"", err, valuesResult)
	}
	if *debugMode {
		println("Debug: rendered ##", renderedValues, "\n---")
	}

	fmt.Println(string(renderedValues))

}
