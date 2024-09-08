package settings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test for GetClientConfig
func TestGetClientConfig(t *testing.T) {
	var config = GetClientConfig()

	// Assertions
	var expected = ClientConfig{
		ServerAddress: ":8080",
		UseTLS:        false,
	}

	assert.Equal(t, expected, config, "Config should match the expected values")
}
