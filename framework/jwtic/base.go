package jwtic

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/dgrijalva/jwt-go"
)

var (
	secretKey = ""

	errInvalidParsedToken = errors.New("parsed token is invalid")
)

// LoadTokenKey ...
func LoadTokenKey(key string) {
	log.Println("set secretKey as: " + key)
	secretKey = key
}

const (
	srcKey = "source"
)

// Encrypt ... from struct into signed string
func Encrypt(v Data) (string, error) {
	if secretKey == "" {
		panic("empty secret key")
	}

	byts, err := v.Marshal()
	if err != nil {
		return "", fmt.Errorf("could not convert struct: %v", err)
	}
	claims := jwt.MapClaims{srcKey: string(byts)}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("could signed v: %v", err)
	}

	return signed, nil
}

// Decrypt ...
func Decrypt(v Data, signed string) error {
	parsedToken, err := jwt.Parse(signed,
		func(t *jwt.Token) (interface{}, error) { return []byte(secretKey), nil })
	if err != nil {
		log.Printf("parsed error: %v\n", err)
		return err
	}

	// log.Printf("parsed token is: %v\n", parsedToken)
	// srcByts, valid := validToken(parsedToken)
	parsedClaims := parsedToken.Claims.(jwt.MapClaims)
	srcByts := parsedClaims[srcKey].(string)

	// if !valid {
	// 	return errInvalidParsedToken
	// }
	return v.Unmarshal([]byte(srcByts))
}

// Data ...
type Data interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}

var (
	_ Data = MapData{}
)

// MapData ...
type MapData map[string]interface{}

// Marshal ...
func (m MapData) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

// Unmarshal ...
func (m MapData) Unmarshal(d []byte) error {
	return json.Unmarshal(d, &m)
}
