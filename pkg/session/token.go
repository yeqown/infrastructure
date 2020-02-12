package session

import (
	"encoding"
	"encoding/json"
	"errors"
	"time"
)

// IToken . an interface which contains necessary information of user who has
// logined to web application
// 支持多平台
type IToken interface {
	// functions for redis encoding and decoding
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler

	// following methods is provided to ITokenManager to op with token
	// TokenKey got key in redis or somewhere to find the token
	TokenKey() string
	// Expiration got time.Duration(int64) of token will expire at
	Expiration() time.Duration
	// Valid will raise error related to the method how to validate the information contained by token
	Valid() error
	// SetExpiration to update token expired time
	SetExpiration(d time.Duration)
}

// DefaultToken . claims for jwt
type DefaultToken struct {
	LastLoginIP   string `json:"last_login_ip"`
	UserID        int64  `json:"user_id"`
	Telephone     string `json:"telepone"`
	Nickname      string `json:"nickname"`
	RedisTokenKey string `json:"token_key"`
	ExpiredAt     int64  `json:"expired_at"` // 时间戳 ms
}

// NewDefaultToken .
func NewDefaultToken(uid int64, tel, nickname, ip string, d time.Duration) IToken {
	us := &DefaultToken{
		LastLoginIP: ip,
		UserID:      uid,
		Telephone:   tel,
		Nickname:    nickname,
		ExpiredAt:   time.Now().Add(d).Unix(),
	}

	us.TokenKey()

	return us
}

// Valid for session .
func (s *DefaultToken) Valid() error {
	if s.UserID == 0 {
		return errors.New("empty userID")
	}

	if s.Telephone == "" {
		return errors.New("empty tel")
	}

	return nil
}

// MarshalBinary .
func (s *DefaultToken) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

// UnmarshalBinary .
func (s *DefaultToken) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}

// TokenKey .
func (s *DefaultToken) TokenKey() string {
	if s.RedisTokenKey == "" {
		s.RedisTokenKey = tokenKey(s.Telephone)
	}
	return s.RedisTokenKey
}

// Expiration .
func (s *DefaultToken) Expiration() time.Duration {
	return time.Duration(s.ExpiredAt-time.Now().Unix()) * time.Millisecond
}

// SetExpiration to update DefaultToken.ExpiredAt timestamp
func (s *DefaultToken) SetExpiration(d time.Duration) {
	s.ExpiredAt = time.Now().Add(d).Unix()
}

func tokenKey(tel string) string {
	return "app:session:" + tel
}
