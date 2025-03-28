package secretstorage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"os"
	"path/filepath"
)

type SecretStorage struct {
	cipher        cipher.Block
	gcm           cipher.AEAD
	sessionIdPath string
	SessionId     *string
}

func New(password string) (*SecretStorage, error) {
	salt := os.Getenv("SECRET_STORAGE_SALT")
	key, err := pbkdf2.Key(sha256.New, password, []byte(salt), 4096, 32)
	if err != nil {
		return nil, err
	}

	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return nil, err
	}

	homePath, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	sessionIdPath := filepath.Join(homePath, ".piximan", "sessionid")

	storage := SecretStorage{aesCipher, gcm, sessionIdPath, nil}

	if _, err := os.Stat(sessionIdPath); err == nil {
		err := storage.retrieveSessionId()
		if err != nil {
			return nil, err
		}
	}

	return &storage, nil
}

func (s *SecretStorage) StoreSessionId(sessionId string) error {
	err := os.MkdirAll(filepath.Dir(s.sessionIdPath), 0775)
	if err != nil {
		return err
	}

	nonce := make([]byte, s.gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return err
	}

	encrypted := s.gcm.Seal(nonce, nonce, []byte(sessionId), nil)

	err = os.WriteFile(s.sessionIdPath, encrypted, 0600)
	if err != nil {
		return err
	}

	s.SessionId = &sessionId
	return nil
}

func (s *SecretStorage) retrieveSessionId() error {
	encrypted, err := os.ReadFile(s.sessionIdPath)
	if err != nil {
		return err
	}

	nonceSize := s.gcm.NonceSize()
	nonce, ciphertext := encrypted[:nonceSize], encrypted[nonceSize:]

	decrypted, err := s.gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return err
	}

	sessionId := string(decrypted)
	s.SessionId = &sessionId
	return nil
}
