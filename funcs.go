// There are some functions like in helm
package main

import (
	"encoding/json"
	"log"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"gopkg.in/yaml.v2"
)

func funcMap() template.FuncMap {
	f := sprig.TxtFuncMap()
	delete(f, "env")
	delete(f, "expandenv")

	extra := template.FuncMap{
		"toYaml": toYAML,
		"toJson": toJSON,

		// functions are not implemented and I don't want to
		"include":  func(string, interface{}) string { return "not implemented" },
		"tpl":      func(string, interface{}) string { return "not implemented" },
		"required": func(string, interface{}) string { return "not implemented" },
	}

	for k, v := range extra {
		f[k] = v
	}

	return f
}

func toYAML(v interface{}) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		log.Fatalf("Error: Can't execute toYAML func:\"%v\"\n   \"%s\"", err, v)
	}
	return strings.TrimSuffix(string(data), "\n")
}

func toJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		log.Fatalf("Error: Can't execute toJSON func:\"%v\"\n   \"%s\"", err, v)
	}
	return string(data)
}
