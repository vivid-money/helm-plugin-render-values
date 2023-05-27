# helm-plugin-render-values

The `helm-plugin-render-values` is a Helm downloader plugin that enhances the rendering of templated values.

This plugin is particularly useful when working with third-party charts that don't have the tpl() function for values, but you still want to utilize it.

Feel free to explore the capabilities of this plugin and enjoy a more streamlined and efficient Helm charting experience!

## Install
Use helm CLI to install this plugin:
```
$ helm plugin install https://github.com/vivid-money/helm-plugin-render-values --version 0.2.2
```

## Usage
```
helm upgrade name . -f render://templated-values.yaml
```
templated-values.yaml should looks like this
```
## Extra files values imorted from. They must not have templating
importValuesFrom: 
- base-values1.yaml
- base-values2.yaml
- services/*/deploy.yaml
- services/*/env/dev.yaml

## Extra files with templating to render. They should have templating
extendRenderWith: 
- extended-values1.yaml

myapp:
  cluster: {{ .Values.clusterName }}
  enabled: {{ .Values.enabled }}
  {{if .Values.default }}
  run: {{ .Values.default }}
  {{- end }}
```
### *importValuesFrom*
- The `importValuesFrom` keyword allows you to specify a list of sources for values to render them.
- These files should not contain any templating. Only values for templating!
- If you use `importValuesFrom` with a pattern, the values will be nested under the corresponding folder structure. For example, if you have values in files under the `folder1/folder2` directory, they will be nested under `"folder1"."folder2"` in the rendered values.
- The values specified in `importValuesFrom` will be merged, and if there are conflicting keys, the last imported value will override the first one.

### *extendRenderWith*
- The `extendRenderWith` section specifies a list of files with templating to render. These files can contain templating logic. In the example, values are extended from extended-values1.yaml

### Self values rendering
If neither `importValuesFrom` nor `extendRenderWith` is specified in the file, the values will be taken from the same file, and the file itself will be used for template rendering.

- Since this file renders values based on its own values, templating will only work with a single level of nesting. For example, if we define param1: value1, we can use it for param2: {{ .Values.param1 }}, but we cannot use param2 to generate the next value.

- The plugin doesn't know about `.Release` values!

- Another limitation is that you cannot use string values with `{{ xxx }}` in the main file that should not be rendered (e.g., passing strings with templating for Prometheus rules).

Example:

`helm upgrade releasename -f render://test-values.yaml`
```test-values.yaml
env: dev
namespace: release-{{.Value.env}}
hostname: service-{{.Value.env}}.domain.com
```
## Notes

v0.2.2
- update go to 1.19
- added tests
- update README

v0.2.1
- importValuesFrom now can been set with files pattern

v0.2.0
- add "extendRenderWith" key to include those files to output render
- fix mergeKeys func

v0.1.3
- fixed all from previous(0.1.2) note
- could use "self" in importValuesFrom

***

## Development

Run tests
`go test ./tests/ -v`
