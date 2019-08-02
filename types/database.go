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

// SQLite3Config .
type SQLite3Config struct {
	Name string // DB file name
}

// MgoConfig .
// mongo host:port/db_name -u username -p password
type MgoConfig struct {
	Addrs     string `json:"addrs"`
	Timeout   int    `json:"timeout"`
	Database  string `json:"database"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	PoolLimit int    `json:"poollimit"`
}

func (c *MgoConfig) String() string {
	return fmt.Sprintf("%v", *c)
}

// RedisConfig .
type RedisConfig struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

func (c *RedisConfig) String() string {
	return fmt.Sprintf("%v", *c)
}
