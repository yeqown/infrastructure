package dbs

/*
TODO:
	a tool to generate Model struct to service Struct
	like convert UserModel to UserService

type UserModel struct {
	Name     string `gorm:"column:name"`
	Password string `gorm:"column:password"`
}

type UserService struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
*/

// 方向：
// 1. go generate
// 2. reflect
