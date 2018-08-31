package tools

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/template"
)

func defaultGenTagFunc(name, typ, tag string) string {
	return fmt.Sprintf("json:\"%s\"", tag)
}

func genStructFieldLine(name, typ, tag string) string {
	return fmt.Sprintf("%s %s `%s`", name, typ,
		genTag(name, typ, tag),
	)
}

func wirteFile(filename string, bs []byte) {
	fd, err := os.OpenFile(filename, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	fd.Write(bs)
	return
}

type outfileCfg struct {
	exportFilename  string
	exportPkgName   string
	modelImportPath string
}

func generateFile(ises []*innerStruct, cfg *outfileCfg) {
	buf := bytes.NewBuffer([]byte{})

	if err := writeFileHeader(buf, cfg.exportPkgName, cfg.modelImportPath); err != nil {
		log.Fatalln(err)
	}

	for _, is := range ises {
		if err := wirteStruct(buf, is); err != nil {
			log.Fatalln(err)
		}

		if err := writeLoadModelFunc(buf, is); err != nil {
			log.Fatalln(err)
		}
	}

	if isDebug {
		log.Println("outfile:", cfg.exportFilename)
	}

	wirteFile(cfg.exportFilename, buf.Bytes())
}

func writeFileHeader(w io.Writer, pkgName, modelImportPath string) error {
	fh := struct {
		PkgName         string
		ModelImportPath string
	}{
		pkgName, modelImportPath,
	}

	tmpl, err := template.New("fh").Parse(fileHeaderTmpl)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, fh)
}

// writerStruct generate struct bytes with structTmpl
func wirteStruct(w io.Writer, is *innerStruct) error {
	ss := struct {
		Name   string
		Fields []string
	}{
		Name: genStructName(is.name),
	}

	// struct fields generating
	for _, fld := range is.fields {
		line := genStructFieldLine(fld.name, fld.typ, fld.tag)
		ss.Fields = append(ss.Fields, line)
	}

	if isDebug {
		log.Println("generate Struct with:", ss)
	}

	tmpl, err := template.New("ss").Parse(structTmpl)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, ss)
}

func writeLoadModelFunc(w io.Writer, is *innerStruct) error {
	f := struct {
		ToStructName    string
		ModelPkgName    string
		ModelStructName string
		Fields          []string
	}{
		ToStructName:    genStructName(is.name),
		ModelPkgName:    is.pkgName,
		ModelStructName: is.name,
	}

	for _, fld := range is.fields {
		f.Fields = append(f.Fields, fld.name)
	}

	tmpl, err := template.New("f").Parse(loadFromModelFuncTmpl)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, f)
}

func genStructName(name string) string {
	name = strings.TrimSuffix(name, specifiedStructTypeSuffix)
	name += exportStructSuffix
	return name
}
