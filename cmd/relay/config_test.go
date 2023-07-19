package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigFromFile(t *testing.T) {
	configFile, err := ioutil.TempFile("", "config.*.yml")
	assert.NoError(t, err)

	defer os.Remove(configFile.Name())

	_, err = configFile.Write([]byte(`
port: 9000
database_url: postgresql://example
nip11_pubkey: fixme
nip11_contact: fixme
`))
	assert.NoError(t, err)

	var cfg Config
	err = cfg.Load(configFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, 9000, cfg.Port)
	assert.Equal(t, "postgresql://example", cfg.DatabaseURL)
	assert.Equal(t, "fixme", cfg.Nip11Pubkey)
	assert.Equal(t, "fixme", cfg.Nip11Contact)
}
