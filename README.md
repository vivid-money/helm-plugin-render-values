# helm-plugin-render-values

The Helm downloader plugin with rendering templated values files

## Install
Use helm CLI to install this plugin:
```
$ helm plugin install https://github.com/vivid-money/helm-render-values --version 0.1.1
```

## Usage
```
upgrade name . -f render://templated-values.yaml
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

