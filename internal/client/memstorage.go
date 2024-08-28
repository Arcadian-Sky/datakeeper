package client

import (
	"log"
	"os"
	"path"
)

// MemStorage stores runtime state of client
type MemStorage struct {
	Login        string
	Token        string
	MasterKey    MasterKey
	MasterKeyDir string
	PfilesDir    string
}

// SetToken sets/updates token
func (m *MemStorage) SetToken(token string) {
	m.Token = token
}

// SetMasterKey sets/updates MasterKey
func (m *MemStorage) SetMasterKey(key string, keyPath string) {
	m.MasterKey.Key = key
	m.MasterKey.KeyPath = keyPath
	m.MasterKey.SetHash()
}

// NewMemStorage returns new MemStorage instance
func NewMemStorage() *MemStorage {
	mstorage := MemStorage{}
	mstorage.MasterKeyDir = createKeyDir()
	mstorage.PfilesDir = createPfileDir()
	return &mstorage
}
func createKeyDir() string {
	userHome, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("cant get user home directory: %s", err.Error())
	}
	kpath := path.Join(userHome, ".gk-keychain")
	err = os.MkdirAll(kpath, 0700)
	if err != nil {
		log.Fatalf("cant create keychain directory(%s): %s", kpath, err.Error())
	}
	return kpath
}

func createPfileDir() string {
	userHome, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("cant get user home directory: %s", err.Error())
	}
	kpath := path.Join(userHome, ".gk-pfiles")
	err = os.MkdirAll(kpath, 0700)
	if err != nil {
		log.Fatalf("cant create pfile directory(%s): %s", kpath, err.Error())
	}
	return kpath
}
