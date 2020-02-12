package session

import (
	"errors"
	"fmt"
	"time"

	logger "github.com/yeqown/infrastructure/framework/logrus-logger"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
)

var (
	// OneDayExpired .
	OneDayExpired = 24 * time.Hour
	// OneWeekExpired .
	OneWeekExpired = 7 * 24 * time.Hour
	// OneMonthExpired .
	OneMonthExpired = 30 * 24 * time.Hour
)

// ITokenManager .
// 1. 能生成sessionkey 和 session 内容
// 2. 并发安全
// 3. 支持多平台token不冲突
type ITokenManager interface {
	Set(tok IToken) error
	Exists(tok IToken) bool
	Expire(tok IToken) error
	Refresh(tok IToken, d time.Duration) error
	Get(key string, tok IToken) error
	PraseToken(token string, tok IToken) error
	Token(tok IToken) (string, error)
}

var (
	errTokenValid = errors.New("invalid token")
)

type defualtTokenManager struct {
	redisClient *redis.Client
	signedKey   []byte
}

// NewJwtTokenManager .
func NewJwtTokenManager(rc *redis.Client, signedKey string) ITokenManager {
	return &defualtTokenManager{
		redisClient: rc,
		signedKey:   []byte(signedKey),
	}
}

// Set .
func (mgr *defualtTokenManager) Set(tok IToken) error {
	key := tok.TokenKey()
	return mgr.redisClient.Set(key, tok, tok.Expiration()).Err()
}

// Exists .
func (mgr *defualtTokenManager) Exists(tok IToken) bool {
	key := tok.TokenKey()
	i, err := mgr.redisClient.Exists(key).Result()
	if err != nil {
		logger.Log.Warnf("defualtTokenManager.Exists(%s) failed, err=%v", key, err)
		return false
	}

	return i == 1
}

// Expire .
func (mgr *defualtTokenManager) Expire(tok IToken) error {
	key := tok.TokenKey()
	return mgr.redisClient.Expire(key, time.Duration(-1)).Err()
}

// Refresh .
func (mgr *defualtTokenManager) Refresh(tok IToken, d time.Duration) error {
	key := tok.TokenKey()
	return mgr.redisClient.Set(key, tok, d).Err()
}

// Get .
func (mgr *defualtTokenManager) Get(key string, tok IToken) error {
	return mgr.redisClient.Get(key).Scan(tok)
}

// PraseToken .
func (mgr *defualtTokenManager) PraseToken(tokString string, tok IToken) error {
	jwtTok, err := jwt.ParseWithClaims(tokString, tok, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return mgr.signedKey, nil
	})

	if err != nil {
		// true: cannot parse token
		return err
	}

	tok, ok := jwtTok.Claims.(IToken)
	if !ok || !jwtTok.Valid {
		// fmt.Println(ok, jwtTok.Valid)
		return errTokenValid
	}

	return nil
}

func (mgr *defualtTokenManager) Token(tok IToken) (string, error) {
	jwtTok := jwt.NewWithClaims(jwt.SigningMethodHS256, tok)
	return jwtTok.SignedString(mgr.signedKey)
}
