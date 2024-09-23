package eventsource

import (
	"testing"
)

func TestSymmetric(t *testing.T) {
	key := "N1PCdw3M2B1TfJhoaY2mL736p2vCUc47"
	text := "plain text"

	s := Symmetric{key: key}

	ciphertext, err := s.Encrypt(text)
	if err != nil {
		t.Fatal(err)
	}

	plaintext, err := s.Decrypt(ciphertext)
	if err != nil {
		t.Fatal(err)
	}

	if plaintext != text {
		t.Fatalf("Plain text does not match original text. Expected: %v, Recieved: %v", text, plaintext)
	}
}
