package render

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

var (
	errLog log.Logger
)

type Values map[interface{}]interface{}

type Files struct {
	ImportValues []string `yaml:"importValuesFrom"`
	ExtendRender []string `yaml:"extendRenderWith"`
}

type ValuesRenderer struct {
	filename string
	files    Files
	values   Values
	Debug    bool
}

func (vr *ValuesRenderer) Debuging(format string, v ...any) {
	if vr.Debug {
		errLog.Printf(format, v...)
	}
}

// Create new render.
func (vr *ValuesRenderer) Run(filename, argPrefix string) {
	errLog.SetOutput(os.Stderr)
	vr.Debuging("DEBUG: filename: %s", string(filename))

	vr.filename = strings.TrimPrefix(filename, argPrefix)

	if err := vr.GetFiles(); err != nil {
		errLog.Fatalf("ERROR: %s", err)
	}
	if err := vr.ReadValues(); err != nil {
		errLog.Fatalf("ERROR:  %s", err)
	}
	if err := vr.RenderTemplate(); err != nil {
		errLog.Fatalf("ERROR:  %s", err)
	}
}

// Get files for importValuesFrom/extendRenderWith.
func (vr *ValuesRenderer) GetFiles() error {

	rawFile := readFile(vr.filename)
	yamlFile := strings.ReplaceAll(string(rawFile), "{{", "#{{")
	extraFiles := Files{}
	dmsgs := []string{}
	err := yaml.Unmarshal([]byte(yamlFile), &extraFiles)
	if err != nil {
		return fmt.Errorf("broken file: \"%s\" reason: \"%v\"", vr.filename, err)
	}
	if len(extraFiles.ImportValues) > 0 {
		for i, importFilename := range extraFiles.ImportValues {
			if importFilename == "self" {
				extraFiles.ImportValues[i] = vr.filename
				dmsgs = append(dmsgs, "importValuesFiles: self file will been used for values\n")
			} else {
				extraFiles.ImportValues[i] = absolutePath(vr.filename, importFilename)
				dmsgs = append(dmsgs, "importValuesFiles: extrafile "+importFilename+" will been used for values\n")
			}
		}
	} else {
		extraFiles.ImportValues = append(extraFiles.ImportValues, vr.filename)
		dmsgs = append(dmsgs, "importValuesFiles: there is no import values. Selfile will been used\n")
	}
	if len(extraFiles.ExtendRender) > 0 {
		for i, renderFilername := range extraFiles.ExtendRender {
			extraFiles.ExtendRender[i] = absolutePath(vr.filename, renderFilername)
			dmsgs = append(dmsgs, "extendRenderWith: extrafile "+renderFilername+" will been rendred\n")
		}
	} else {
		dmsgs = append(dmsgs, "extendRenderWith: there is no extended files for render\n")
	}

	extraFiles.ExtendRender = append(extraFiles.ExtendRender, vr.filename)
	vr.files = extraFiles

	vr.Debuging("DEBUG: %v\n", dmsgs)
	return err
}

// Read values from files.
func (vr *ValuesRenderer) ReadValues() error {

	vals := make(Values)

	for _, file := range vr.files.ImportValues {
		var rawFile []byte
		var data Values

		if strings.Contains(file, "*") {
			println("Read files with mask")
			data = ParseYamlGlogFile(file)
		} else {
			rawFile = readFile(file)
			if file == vr.filename {
				yamlFile := strings.ReplaceAll(string(rawFile), "{{", "#{{")
				rawFile = []byte(yamlFile)
			}
			err := yaml.Unmarshal(rawFile, &data)
			if err != nil {
				return fmt.Errorf("broken file: \"%s\" reason: \"%v\"", file, err)
			}
		}
		mergeKeys(vals, data)
	}
	if len(vals) == 0 {
		vals = Values{}
	}
	vr.values = make(Values)

	vr.Debuging("DEBUG: total values: %#v\n", vals)
	vr.values["Values"] = vals
	return nil
}

// Read glob files.
func ParseYamlGlogFile(pattern string) Values {

	var data Values
	matches, err := filepath.Glob(pattern)
	if err != nil {
		errLog.Fatalf("ERROR: wrong glob: \"%s\"; stack:\"%v\"", pattern, err)
	}
	for _, file := range matches {

		dir := strings.Split(filepath.Dir(file), "/")
		fileValue := make(Values)

		rawFile := readFile(file)
		err := yaml.Unmarshal(rawFile, &fileValue)
		if err != nil {
			errLog.Fatalf("ERROR: ReadValues: Can't parse file: \"%s\"; stack:\"%v\"", file, err)
		}
		mergeKeys(data, DirrectoryMapping(dir, fileValue))
	}
	return data
}

// Render Values template to a stdout.
func (vr *ValuesRenderer) RenderTemplate() error {

	valuesResult := make(Values)
	for _, file := range vr.files.ExtendRender {
		var data Values
		var buf strings.Builder
		tpl, err := template.New("render").Funcs(funcMap()).ParseFiles(file)

		if err != nil {
			return fmt.Errorf("can't create template: %s}", err)
		}
		tpl.Option("missingkey=error")

		err = tpl.ExecuteTemplate(&buf, filepath.Base(file), vr.values)
		if err != nil {
			return fmt.Errorf("can't render template: %s}", err)
		}
		rendered := strings.ReplaceAll(buf.String(), "<no value>", "")

		// merge output in a single YAML
		err = yaml.Unmarshal([]byte(rendered), &data)
		if err != nil {
			return fmt.Errorf("can't  parse file: \"%s\"; stack:\"%s\"", file, err)
		}
		mergeKeys(valuesResult, data)
	}
	renderedValues, err := yaml.Marshal(valuesResult)
	vr.Debuging("DEBUG: rendered: %#v", string(renderedValues))
	if err != nil {
		return fmt.Errorf("can't marshal yaml: \"%#v\"; stack:\"%v\"", valuesResult, err)
	}

	fmt.Println(string(renderedValues))
	return nil

}
