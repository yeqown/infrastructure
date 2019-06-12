package lang

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// LoadJSONConfigFile filename string
func LoadJSONConfigFile(filename string, v interface{}) error {
	fd, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(fd)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, v)
}

// LoadYamlConfigFile need ?
func LoadYamlConfigFile(filename string, v interface{}) error {
	return nil
}
