package client

import (
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMasterKey_SetHash(t *testing.T) {
	key := "testKey"
	expectedHash := sha256.Sum256([]byte(key))

	masterKey := &MasterKey{Key: key}
	masterKey.SetHash()

	// Verify that KeyHash is correctly set
	assert.Equal(t, expectedHash[:], masterKey.KeyHash, "Expected KeyHash to be the SHA256 hash of Key")
}

func TestMasterKey_Str(t *testing.T) {
	key := "testKey"
	keyPath := "/path/to/key"
	hash := sha256.Sum256([]byte(key))

	masterKey := &MasterKey{
		Key:     key,
		KeyPath: keyPath,
		KeyHash: hash[:],
	}

	expectedStr := fmt.Sprintf("<MasterKey key:'%s', keyPath:'%s' keyHash: '%x'>", key, keyPath, hash[:])
	actualStr := masterKey.Str()

	// Verify that Str returns the expected string representation
	assert.Equal(t, expectedStr, actualStr, "Expected Str to return the formatted string representation of MasterKey")
}
