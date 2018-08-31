package tools

// 解析一个go文件获取其中之类特征的struct并解析结构
import (
	"go/types"
	"log"
	"path"
	"strings"

	"golang.org/x/tools/go/loader"
)

type parseTagFunc func(s string) string

var (
	conf                      loader.Config
	parseTag                  parseTagFunc = defaultParseTagFunc
	specifiedStructTypeSuffix              = "Model"
	specifiedStructTypePrefix              = "Model"
	isDebug                   bool
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

// Exported, and specified type
func loadGoFiles(dir string, importPath string, filenames ...string) ([]*innerStruct, error) {
	conf.Cwd = dir
	conf.CreateFromFilenames(importPath, filenames...)

	prog, err := conf.Load()
	if err != nil {
		log.Println("load program err:", err)
		return nil, err
	}

	if isDebug {
		log.Println("dir:", dir)
		log.Println("importPath:", importPath)
		log.Println("filename:", filenames)
		// log.Println("len of prog imported", len(prog.Imported))
	}

	return loopProgramCreated(prog.Created), nil
}

// loopProgramCreated to loo and filter:
// 1. unexported type
// 2. bultin types
// 3. only specified style struct name
func loopProgramCreated(
	created []*loader.PackageInfo,
) (innerStructs []*innerStruct) {
	for _, pkgInfo := range created {
		pkgName := pkgInfo.Pkg.Name()
		defs := pkgInfo.Defs

		// imports := pkgInfo.Pkg.Imports()
		// for _, imp := range imports {
		// 	log.Println(imp.Path(), imp.Name())
		// }

		for indent, obj := range defs {
			if !indent.IsExported() ||
				obj == nil ||
				!strings.HasSuffix(indent.Name, specifiedStructTypeSuffix) {
				continue
			}

			st, ok := obj.Type().Underlying().(*types.Struct)
			if !ok {
				log.Println("not a struct, skip this")
				continue
			}
			is := new(innerStruct)

			is.content = st.String()
			is.pkgName = pkgName
			is.name = obj.Name()
			is.fields = parseStructFields(st)

			if isDebug {
				log.Println("parse one Model: ", is.name, is.pkgName, is.content)
			}

			innerStructs = append(innerStructs, is)
		}
	}
	return
}

type field struct {
	name string
	typ  string
	tag  string
}

type innerStruct struct {
	fields  []*field
	content string
	name    string
	pkgName string
}

// parseStructFields parse fields
func parseStructFields(st *types.Struct) []*field {
	flds := make([]*field, 0, st.NumFields())

	for i := 0; i < st.NumFields(); i++ {
		fld := st.Field(i)
		isField := new(field)

		isField.name = fld.Name()
		isField.tag = parseTag(st.Tag(i))
		isField.typ = fld.Type().String()

		flds = append(flds, isField)
	}
	return flds
}

/*
func (is *innerStruct) pureName() {
	is.name = strings.TrimPrefix(is.name, is.pkgName+".")
}

// s usually like this:
// type testdata.UserModel struct{Name string "gorm:\"colunm:name\""; Password string "gorm:\"column:password\""}
func parseStructString(s string) *innerStruct {
	s = strings.TrimSpace(s)
	splited := strings.Split(s, " ")

	if len(splited) < 3 {
		log.Fatalf("parseStructString error: input %s, output: %v\n", s, splited)
	}
	content := strings.Join(splited[2:], " ")
	content = strings.TrimPrefix(content, "struct")

	return &innerStruct{
		fields:  parseStructFields(content),
		content: content,
		name:    splited[1],
	}

}


// s usually like this:
// {Name string "gorm:\"colunm:name\""; Password string "gorm:\"column:password\""}
func parseStructFields(s string) []*field {
	// s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "{")
	s = strings.TrimSuffix(s, "}")

	splited := strings.Split(s, ";")
	fields := make([]*field, 0, len(splited))
	for _, fldStr := range splited {
		fldStr = strings.TrimSpace(fldStr)
		if fldStr == "" {
			continue
		}
		if isDebug {
			log.Println("parsing field string: ", fldStr)
		}
		fields = append(fields, parseField(fldStr))
	}
	return fields
}

// s usually like this:
// Name string "gorm:\"colunm:name\""
func parseField(s string) *field {
	s = strings.TrimSpace(s)
	splited := strings.Split(s, " ")

	if len(splited) < 3 {
		log.Fatalf("parseFiled error: input %s, output: %v\n", s, splited)
	}

	// log.Println(splited)

	tag := cleanTag(splited[2])
	tag = parseTag(tag)

	if isDebug {
		log.Println("parseTag result: ", tag)
	}

	return &field{
		name: splited[0],
		typ:  splited[1],
		tag:  tag,
	}
}

// input: "gorm:\"colunm:name\""
// output: gorm:"column:name"
func cleanTag(tag string) string {
	tag = strings.Replace(tag, `\`, "", -1)
	tag = strings.TrimSuffix(tag, `"`)
	tag = strings.TrimPrefix(tag, `"`)
	return tag
}
*/
// input: gorm:"colunm:name"
// output: gorm:column:name
func defaultParseTagFunc(s string) string {
	s = strings.Replace(s, `"`, "", -1)
	splited := strings.Split(s, ":")
	return splited[len(splited)-1]
}
