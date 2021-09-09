# helm-plugin-render-values

The Helm downloader plugin with rendering templated values files

## Install
Use helm CLI to install this plugin:
```
$ helm plugin install https://github.com/vivid-money/helm-plugin-render-values --version 0.1.3
```

## Usage
```
helm upgrade name . -f render://templated-values.yaml
```
templated-values.yaml should looks like this
```
importValuesFrom: 
- base-values1.yaml
- base-values2.yaml

myapp:
  cluster: {{ .Values.clusterName }}
  enabled: {{ .Values.enabled }}
```

**importValuesFrom** - is a keyword for list with sources for Values to render it

## Notes

v0.1.3
- fixed all from previous(0.1.2) note
- could use "self" in importValuesFrom

v0.1.2 (not actual)
- go template [Actions](https://pkg.go.dev/text/template#hdr-Actions) could be used but yaml should we readable for yaml-parsers. So easest way to do it is using comment befor actions "#" 
like this:
```
myapp:
# {{ if .Values.istio }}
  virtualService: {{ .Values.hostname }}
# {{ else }}
  ingress: {{ .Values.hostname }}
# {{end }}
```

- Don't use helm function for trim sapces "{{-" or "-}}" - it isn't implemented
