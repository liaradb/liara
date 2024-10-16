package shred

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

type Symmetric struct {
	key string
	gcm cipher.AEAD
}

func NewSymmetric(key string) *Symmetric {
	return &Symmetric{key: key}
}

func (s *Symmetric) init() error {
	if s.gcm != nil {
		return nil
	}

	aes, err := aes.NewCipher([]byte(s.key))
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return err
	}

	s.gcm = gcm

	return nil
}

func (s *Symmetric) nonce() ([]byte, error) {
	nonce := make([]byte, s.gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	return nonce, nil
}

func (s *Symmetric) Encrypt(plaintext string) (string, error) {
	err := s.init()
	if err != nil {
		return "", err
	}

	nonce, err := s.nonce()
	if err != nil {
		return "", err
	}

	ciphertext := s.gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return string(ciphertext), nil
}

func (s *Symmetric) Decrypt(ciphertext string) (string, error) {
	err := s.init()
	if err != nil {
		return "", err
	}

	nonceSize := s.gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := s.gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return "", nil
	}

	return string(plaintext), nil
}
