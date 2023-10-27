package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseGriffonBlock(t *testing.T) {
	src := []byte(`
	griffon {
		region = "nyc1"
		vultr_api_key = "1234567890"
	}`)

	config, err := ParseHCL("test.hcl", src)
	require.NoError(t, err)
	require.Equalf(t, GriffonBlock{
		Region:      "nyc1",
		VultrAPIKey: "1234567890",
	}, config.Griffon, "GriffonBlock parsed incorrectly")
}
