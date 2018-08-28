package tools

// 解析一个go文件获取其中之类特征的struct并解析结构
import (
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
	ExportDir          string
	ExportFilename     string
	ExportPkgName      string
	ExportStructSuffix string
	ModelImportPath    string
	StructSuffix       string
	Debug              bool
}

// ParseAndGenerate parse all input go files and
// get wanted struct info save with innerStruct, and then generate file
func ParseAndGenerate(uCfg *UsageCfg, dir string, filenames ...string) error {
	specifiedStructTypeSuffix = uCfg.StructSuffix
	isDebug = uCfg.Debug

	// parse
	ises, err := loadGoFiles(dir, filenames...)
	if err != nil {
		return err
	}

	if isDebug {
		log.Println("loadGoFiles got Struct count:", len(ises))
	}

	// generate
	generateFile(ises, &outfileCfg{
		exportFilename:  path.Join(uCfg.ExportDir, uCfg.ExportFilename),
		exportPkgName:   uCfg.ExportPkgName,
		modelImportPath: uCfg.ModelImportPath,
	})

	return nil
}

// SetCustomParseTagFunc use user's custom parseTag func
func SetCustomParseTagFunc(f parseTagFunc) {
	parseTag = f
}

// Exported, and specified type
func loadGoFiles(dir string, filenames ...string) ([]*innerStruct, error) {
	newFilenames := []string{}

	for _, filename := range filenames {
		newFilenames = append(newFilenames, path.Join(dir, filename))
	}

	conf.CreateFromFilenames("", newFilenames...)
	prog, err := conf.Load()
	if err != nil {
		log.Println("load program err:", err)
		return nil, err
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

		// log.Println(pkgName)
		// for astExpr, typAndVal := range typs {
		// 	if !typAndVal.IsType() && typAndVal.IsBuiltin() {
		// 		continue
		// 	}
		// 	log.Println(typAndVal, astExpr)
		// }

		for indent, obj := range defs {

			if !indent.IsExported() ||
				obj == nil ||
				!strings.HasSuffix(indent.Name, specifiedStructTypeSuffix) {
				continue
			}

			// log.Println(indent, obj)
			// log.Println(indent.Name, indent.String(), obj.String())
			// log.Println(obj)
			is := parseStructString(obj.String())
			is.pkgName = pkgName
			is.pureName()

			if isDebug {
				log.Println("parse one Model: ", is.name, is.pkgName, is.content)
			}

			innerStructs = append(innerStructs, is)
		}
	}
	return
}

type innerStruct struct {
	fields  []*field
	content string
	name    string
	pkgName string
}

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

type field struct {
	name string
	typ  string
	tag  string
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

// input: gorm:"colunm:name"
// output: gorm:column:name
func defaultParseTagFunc(s string) string {
	s = strings.Replace(s, `"`, "", -1)
	splited := strings.Split(s, ":")
	return splited[len(splited)-1]
}
