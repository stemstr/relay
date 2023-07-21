package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port             int               `yaml:"port"`
	DatabaseURL      string            `yaml:"database_url"`
	AllowedKinds     []int             `yaml:"allowed_kinds"`
	AllowedPubkeys   []string          `yaml:"allowed_pubkeys"`
	Nip11Pubkey      string            `yaml:"nip11_pubkey"`
	Nip11Contact     string            `yaml:"nip11_contact"`
	Nip11Description string            `yaml:"nip11_description"`
	Nip11Version     string            `yaml:"nip11_version"`
	Admins           map[string]string `yaml:"admins"`
}

// Load Config from a yaml file at path.
func (c *Config) Load(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return yaml.NewDecoder(f).Decode(c)
}
