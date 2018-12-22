package core

import (
	"io"

	"gopkg.in/yaml.v2"
)

type YAMLConfigParser struct{}

func (y *YAMLConfigParser) Decode(r io.Reader) (*Config, error) {
	decoder := yaml.NewDecoder(r)
	c := new(Config)
	err := decoder.Decode(c)
	return c, err
}

func (y *YAMLConfigParser) Encode(w io.Writer, c *Config) error {
	encoder := yaml.NewEncoder(w)
	defer encoder.Close()
	return encoder.Encode(c)
}
