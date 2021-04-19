package config

import (
	"encoding/json"
	"io"
)

// JSONConfigParser implements the Parser interface
// for a JSON config file.
type JSONConfigParser struct{}

func (y *JSONConfigParser) Decode(r io.Reader) (*Config, error) {
	decoder := json.NewDecoder(r)
	c := new(Config)
	err := decoder.Decode(c)
	return c, err
}

func (y *JSONConfigParser) Encode(w io.Writer, c *Config) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(c)
}
