package cfgutil

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
)

var (
	// ErrInvalidJSONFormat .
	ErrInvalidJSONFormat = errors.New("invalid json format")
)

// LoadJSON .
func LoadJSON(reader io.Reader, recv interface{}) error {
	byts, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	if !json.Valid(byts) {
		return ErrInvalidJSONFormat
	}

	return json.Unmarshal(byts, recv)
}

// Open .
func Open(fp string) (io.ReadCloser, error) {
	return os.OpenFile(fp, os.O_RDONLY, 0644)
}
