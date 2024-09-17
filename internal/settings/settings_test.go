package settings

import (
	"encoding/hex"
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGenerateSecretKey tests the GenerateSecretKey function
func TestGenerateSecretKey(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name         string
		length       int
		expectErr    bool
		expectLength int
	}{
		{"Valid length 16", 16, false, 32}, // 16 bytes = 32 hex characters
		{"Valid length 32", 32, false, 64}, // 32 bytes = 64 hex characters
		{"Zero length", 0, false, 0},       // 0 bytes = 0 hex characters
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			key, err := GenerateSecretKey(tc.length)

			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, key, tc.expectLength)
				// Check if the key is a valid hex string
				_, decodeErr := hex.DecodeString(key)
				assert.NoError(t, decodeErr)
			}
		})
	}
}

func TestParse_DefaultValues(t *testing.T) {
	// Reset environment variables
	t.Setenv("DATAKEEPER_RUN_ADDRESS", "")
	t.Setenv("PG_DATABASE_URI", "")
	t.Setenv("MG_DATABASE_URI", "")
	t.Setenv("FILE_DATABASE_URI", "")
	t.Setenv("FILE_DATABASE_ACCESS_KEY", "")
	t.Setenv("FILE_DATABASE_SECRET", "")

	// Reset flags
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Run Parse function
	flags := Parse()

	// Validate the results
	assert.Equal(t, "localhost:8080", flags.Endpoint)
	assert.Equal(t, "", flags.DBPGSettings)
	assert.Equal(t, "", flags.DBMGSettings)
	assert.NotEmpty(t, flags.SecretKey) // Ensure secret key is generated
	assert.Equal(t, Storage{
		Endpoint:    "",
		AccessKeyID: "",
		Secret:      "",
	}, flags.Storage)
}

func TestParse_WithEnvironmentVariables(t *testing.T) {
	// Set environment variables
	t.Setenv("DATAKEEPER_RUN_ADDRESS", "http://localhost:9090")
	t.Setenv("PG_DATABASE_URI", "postgres://user:password@localhost/db")
	t.Setenv("MG_DATABASE_URI", "mongodb://user:password@localhost/db")
	t.Setenv("FILE_DATABASE_URI", "http://file.storage")
	t.Setenv("FILE_DATABASE_ACCESS_KEY", "access-key-id")
	t.Setenv("FILE_DATABASE_SECRET", "secret")

	// Reset flags
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Run Parse function
	flags := Parse()

	// Validate the results
	assert.Equal(t, "http://localhost:9090", flags.Endpoint)
	assert.Equal(t, "postgres://user:password@localhost/db", flags.DBPGSettings)
	assert.Equal(t, "mongodb://user:password@localhost/db", flags.DBMGSettings)
	assert.NotEmpty(t, flags.SecretKey) // Ensure secret key is generated
	assert.Equal(t, Storage{
		Endpoint:    "http://file.storage",
		AccessKeyID: "access-key-id",
		Secret:      "secret",
	}, flags.Storage)
}
