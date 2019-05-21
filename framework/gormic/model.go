package gormic

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"hash"
	"io"
	"time"
)

var (
	md5h           hash.Hash
	buf, out, out2 []byte
)

func init() {
	md5h = md5.New()
	buf = make([]byte, 48)
	// n1 := base64.StdEncoding.EncodedLen(48) = 64
	out = make([]byte, 64)
	// n2 := hex.EncodedLen(len(out)) = 128
	out2 = make([]byte, 128)
	// log.Println(n1, n2)
}

func md55(dst, src []byte) {
	md5h.Write(src)
	hex.Encode(dst, md5h.Sum(nil))
}

func newID() string {
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		return ""
	}

	base64.StdEncoding.Encode(out, buf)
	md55(out2, out)

	return string(out2)
}

// Model ... as custom gorm model with some callbacks
// BeforeCreate ... create 之前调用
// BeforeUpdate ... update 之前调用
type Model struct {
	ID         string    `gorm:"column:id;index;not null;type:varchar(128)"`
	CreateTime time.Time `gorm:"column:create_time"`
	UpdateTime time.Time `gorm:"column:update_time"`
}

// BeforeCreate ....
func (m *Model) BeforeCreate() error {
	if m.ID == "" {
		m.ID = newID()
	}
	m.CreateTime = time.Now()
	m.UpdateTime = time.Now()
	return nil
}

// BeforeUpdate ....
func (m *Model) BeforeUpdate() error {
	m.UpdateTime = time.Now()
	return nil
}
