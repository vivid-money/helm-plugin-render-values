# helm-plugin-render-values

The Helm downloader plugin with rendering templated values files

## Install
Use helm CLI to install this plugin:
```
$ helm plugin install https://github.com/vivid-money/helm-plugin-render-values --version 0.2.0
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
- services/*/deploy.yaml
- services/*/env/dev.yaml

extendRenderWith:
- extended-values1.yaml

myapp:
  cluster: {{ .Values.clusterName }}
  enabled: {{ .Values.enabled }}
```

if you use importValuesFrom with a pattern, then the value will be nested in "folder1"."folder2": (values in files)

the values in importValuesFrom will be merged
and the last one will override the first one if they have the same key

look an example in a "test" dir

**importValuesFrom** - is a keyword for list with sources for Values to render it

## Notes

v0.2.1
- importValuesFrom now can been set with files pattern

v0.2.0
- add "extendRenderWith" key to include those files to output render
- fix mergeKeys func

v0.1.3
- fixed all from previous(0.1.2) note
- could use "self" in importValuesFrom

v0.1.2 (not actual)
***
