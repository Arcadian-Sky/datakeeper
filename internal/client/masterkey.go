package client

import (
	"crypto/sha256"
	"fmt"
)

// MasterKey represents master key
type MasterKey struct {
	Key     string
	KeyPath string
	KeyHash []byte
}

// SetHash calculates key hash and stores it in KeyHash attribute
func (m *MasterKey) SetHash() {
	hash := sha256.Sum256([]byte(m.Key))
	m.KeyHash = hash[:]
}

// Str returns stirng representation of struct
func (m *MasterKey) Str() string {
	return fmt.Sprintf("<MasterKey key:'%s', keyPath:'%s' keyHash: '%x'>",
		m.Key, m.KeyPath, m.KeyHash)
}
