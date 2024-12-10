package utils

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

func PublicKeyFile(file string) (ssh.AuthMethod, error) {
	buffer, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("Failed to read identity file %s: %w", file, err)
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse private key from file %s: %w", file, err)
	}

	return ssh.PublicKeys(key), nil
}
