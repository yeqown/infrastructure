package types

import "fmt"

// MysqlC mysql server config struct
type MysqlC struct {
	UserName  string `json:"username"`
	Password  string `json:"password"`
	Addr      string `json:"addr"`
	DBName    string `json:"dbname"`
	Charset   string `json:"charset"`
	ParseTime bool   `json:"parsetime"`
	Loc       string `json:"loc"`
	Pool      int    `json:"pool"`
}

func (c *MysqlC) valid() bool {
	return c.UserName == "" || c.Password == "" || c.DBName == "" ||
		c.Charset == "" || c.Loc == "" || c.Addr == ""
}

// @output "user:password@addr:port/dbname?charset=utf8&parseTime=True&loc=Local"
func (c *MysqlC) String() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%v&loc=%s", c.UserName, c.Password, c.Addr,
		c.DBName, c.Charset, c.ParseTime, c.Loc)
}
