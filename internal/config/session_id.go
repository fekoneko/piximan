package config

import (
	"crypto/rand"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

func (s *Config) SessionId() (*string, error) {
	if s.sessionId != nil {
		return s.sessionId, nil
	}

	encrypted, err := os.ReadFile(sessionIdPath)
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	nonceSize := s.gcm.NonceSize()
	nonce, ciphertext := encrypted[:nonceSize], encrypted[nonceSize:]

	decrypted, err := s.gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return nil, err
	}

	sessionId := string(decrypted)
	s.sessionId = &sessionId

	return s.sessionId, nil
}

func (s *Config) WriteSessionId(sessionId string) error {
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

	s.sessionId = &sessionId
	return nil
}

func (s *Config) ResetSessionId() error {
	err := os.Remove(sessionIdPath)
	if err == nil {
		s.sessionId = nil
	}
	return err
}
