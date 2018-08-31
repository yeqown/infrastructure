package dbs

### 目标:
a tool to generate Model struct to service Struct like convert UserModel to UserService

```Golang
type UserModel struct {
	Name     string `gorm:"column:name"`
	Password string `gorm:"column:password"`
}

type UserService struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
```

### 使用方法

```bash
➜  server-common git:(master) ✗ ./dto -h
Usage of ./dto:
  -debug
        debug mode, if open this will output info
  -dir string
        model directory from where
  -filename string
        specified filename those you want to generate
  -generateDir string
        generate file will be saved here (default ".")
  -generateFilename string
        generate file name will be use this (default "types.go")
  -generatePkgName string
        generatePkgName (default "types")
  -generateStructSuffix string
        replace model struct name suffix, like: UserSuffix => User
  -modelImportPath my-server/models
        model package path, cannot be empty, like my-server/models
  -modelStructSuffix string
        specified in which Model name style can be generate (default "Model")
```

#### 安装

```sh
go get github.com/yeqown/server-common/dbs/tools
# 获取 tool.main.go, 并选择性的实现自己的 CustomParseTagFunc & CustomGenerateTagFunc
go build -o dto tool.main.go
```


#### 使用实例

```sh
# 根据./dbs/tools/testdata目录下的type_model.go生成文件，
# 在./dbs/tools/testdata目录下生成文件，
# 指定包名为testdata
dto -dir=./dbs/tools/testdata -filename=type_model.go -generatePkgName=testdata 
-generateDir=./dbs/tools/testdata -modelImportPath=model

# 多文件可以使用多次 -filename=model1.go -filename=model2.go
```

```golang
// type_model.go
package testdata

import (
	"errors"
	"time"
)

type UserModel struct {
	Name       string    `gorm:"colunm:name"`
	Password   string    `gorm:"column:password"`
	CreateTime time.Time `gorm:"colunm:create_time"`
	UpdateTime time.Time `gorm:"colunm:update_time"`
}

type userStruct struct {
	Name     string `json:"typeStructName"`
	Password string `json:"typeStructPassword"`
}

type AliasString string

var (
	A UserModel
	B int64
	C string
)

func f(a string) error {
	return errors.New(a)
}

// 生成的文件 types.go[未被格式化]

// Package testdata ...
// Generate by github.com/yeqown/server-common/dbs/tools
package testdata

import (
	"model"
)

// User description here
type User struct {
	Name string `json:"name"`
	Password string `json:"password"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
	}
// LoadUserFromModel func to load data from model
func LoadUserFromModel(data *testdata.UserModel) *User {
	return &User {
		Name: data.Name,
		Password: data.Password,
		CreateTime: data.CreateTime,
		UpdateTime: data.UpdateTime,
		}
}
```

#### 配置自定义Tag解析和生成函数

在[tool.main.go](#)文件中已经定义了

```golang
func main() {
	// set custom funcs
	tools.SetCustomGenTagFunc(CustomGenerateTagFunc)
	tools.SetCustomParseTagFunc(CustomParseTagFunc)
}

// CustomParseTagFunc to custom implment yourself parseTagFunc
// @param (gorm:"colunm:name")
// return ("name")
func CustomParseTagFunc(s string) string {
	log.Println("calling  CustomParseTagFunc", s)

	s = strings.Replace(s, `"`, "", -1)
	splited := strings.Split(s, ":")
	return splited[len(splited)-1]
}

// CustomGenerateTagFunc to implment yourself generateTagFunc
// @param name fieldName (Age, Name, Year, CreateTime)
// @param typ fieldType (string, int64, time.Time)
// @param tag (CustomParseTagFunc) return value, default is gorm tag
// return (json:"name")
func CustomGenerateTagFunc(name, typ, tag string) string {
	return fmt.Sprintf("json:\"%s\"", tag)
}

```
