package tools

import (
	"path"
)

type parseTagFunc func(s string) string
type genTagFunc func(name, typ, tag string) string

var (
	genTag                    genTagFunc   = defaultGenTagFunc
	parseTag                  parseTagFunc = defaultParseTagFunc
	exportStructSuffix        string
	specifiedStructTypeSuffix = "Model"
	specifiedStructTypePrefix = "Model"

	isDebug bool
)

// UsageCfg ... config tools some feature
type UsageCfg struct {
	// Dir 设置需要对那一路径下的文件进行解析
	Dir string
	// ExportDir 设置生成的文件的存放地址
	ExportDir string
	// ExportFilename 指定生成的文件名字
	ExportFilename string
	// ExportPkgName 指定生成的文件的包名
	ExportPkgName string
	// ExportStructSuffix 指定生成的新的机构体后缀
	ExportStructSuffix string
	// ModelImportPath 指定源文件所在包的导入路径
	ModelImportPath string
	// StructSuffix 需要解析的自定义结构体后缀
	StructSuffix string
	// Debug 调试模式开关
	Debug bool
	// Filenames 指定需要解析的.go源文件 文件名字
	Filenames []string
}

// ParseAndGenerate parse all input go files and
// get wanted struct info save with innerStruct, and then generate file
func ParseAndGenerate(cfg *UsageCfg) error {
	specifiedStructTypeSuffix = cfg.StructSuffix
	isDebug = cfg.Debug

	// parse
	ises, err := loadGoFiles(cfg.Dir, cfg.ModelImportPath, cfg.Filenames...)
	if err != nil {
		return err
	}

	// generate
	generateFile(ises, &outfileCfg{
		exportFilename:  path.Join(cfg.ExportDir, cfg.ExportFilename),
		exportPkgName:   cfg.ExportPkgName,
		modelImportPath: cfg.ModelImportPath,
	})

	return nil
}

// SetCustomParseTagFunc use user's custom parseTag func
func SetCustomParseTagFunc(f parseTagFunc) {
	parseTag = f
}

// SetCustomGenTagFunc use user's custom genTagFunc
func SetCustomGenTagFunc(f genTagFunc) {
	genTag = f
}
