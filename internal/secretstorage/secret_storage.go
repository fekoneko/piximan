package secretstorage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

var homePath, _ = os.UserHomeDir()
var sessionIdPath = filepath.Join(homePath, ".piximan", "sessionid")

type SecretStorage struct {
	cipher    cipher.Block
	gcm       cipher.AEAD
	SessionId *string
}

func Open(password string) (*SecretStorage, error) {
	// TODO: maybe make the salt not empty and store it as well
	key, err := pbkdf2.Key(sha256.New, password, []byte{}, 4096, 32)
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

	return &SecretStorage{aesCipher, gcm, nil}, nil
}

func (s *SecretStorage) Read() error {
	if _, err := os.Stat(sessionIdPath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		} else {
			return err
		}
	}

	encrypted, err := os.ReadFile(sessionIdPath)
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

func (s *SecretStorage) WriteSessionId(sessionId string) error {
	err := os.MkdirAll(filepath.Dir(sessionIdPath), 0775)
	if err != nil {
		return err
	}

	nonce := make([]byte, s.gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return err
	}

	encrypted := s.gcm.Seal(nonce, nonce, []byte(sessionId), nil)

	err = os.WriteFile(sessionIdPath, encrypted, 0600)
	if err != nil {
		return err
	}

	s.SessionId = &sessionId
	return nil
}

func RemoveSessionId() error {
	if err := os.Remove(sessionIdPath); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
	return nil
}
