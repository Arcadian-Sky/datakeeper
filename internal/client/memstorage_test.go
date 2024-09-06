package client

import (
	"crypto/sha256"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemStorage_SetToken(t *testing.T) {
	storage := &MemStorage{}
	expectedToken := "testToken"

	storage.SetToken(expectedToken)

	// Verify that the Token field is correctly set
	assert.Equal(t, expectedToken, storage.Token, "Expected Token to be set correctly")
}

func TestMemStorage_SetMasterKey(t *testing.T) {
	storage := &MemStorage{}
	key := "testKey"
	keyPath := "/path/to/key"
	expectedHash := sha256.Sum256([]byte(key))

	storage.SetMasterKey(key, keyPath)

	// Verify that the MasterKey fields are correctly set
	assert.Equal(t, key, storage.MasterKey.Key, "Expected MasterKey.Key to be set correctly")
	assert.Equal(t, keyPath, storage.MasterKey.KeyPath, "Expected MasterKey.KeyPath to be set correctly")
	assert.Equal(t, expectedHash[:], storage.MasterKey.KeyHash, "Expected MasterKey.KeyHash to be set correctly")
}
