package tools

// define some tmplates to output format

const fileHeaderTmpl = `// Package {{.PkgName}} ...
// Generate by github.com/yeqown/infrastructure/framework/gormic/tools
package {{.PkgName}}
{{if ne .ModelImportPath ""}}
import (
	"{{.ModelImportPath}}"
)
{{end}}
`
const structTmpl = `// {{.Name}} description here
type {{.Name}} struct {
	{{range $index, $fieldline := .Fields}}{{$fieldline}}
	{{end}}}
`

const loadFromModelFuncTmpl = `// Load{{.ToStructName}}FromModel func to load data from model
func Load{{.ToStructName}}FromModel(data *{{.ModelPkgName}}.{{.ModelStructName}}) *{{.ToStructName}} {
	return &{{.ToStructName}} {
		{{range $index, $fld := .Fields}}{{$fld}}: data.{{$fld}},
		{{end}}}
}
`
