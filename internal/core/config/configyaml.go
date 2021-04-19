package config

import (
	"io"

	"github.com/ghodss/yaml"
)

// YAMLConfigParser implements the Parser interface
// for a YAML config file.
type YAMLConfigParser struct{}

func (y *YAMLConfigParser) Decode(r io.Reader) (*Config, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	c := new(Config)
	err = yaml.Unmarshal(data, c)
	return c, err
}

func (y *YAMLConfigParser) Encode(w io.Writer, c *Config) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}
