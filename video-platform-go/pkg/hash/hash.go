package hash

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

type HashConfig struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

var DefaultConfig = HashConfig{
	Memory:      64 * 1024,
	Iterations:  3,
	Parallelism: 2,
	SaltLength:  16,
	KeyLength:   32,
}

func HashPassword(password string) (string, error) {
	return HashPasswordWithConfig(password, DefaultConfig)
}

func HashPasswordWithConfig(password string, config HashConfig) (string, error) {
	salt := make([]byte, config.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		config.Iterations,
		config.Memory,
		config.Parallelism,
		config.KeyLength,
	)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		config.Memory,
		config.Iterations,
		config.Parallelism,
		b64Salt,
		b64Hash,
	)

	return encodedHash, nil
}

func CheckPassword(password, encodedHash string) bool {
	config, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false
	}

	otherHash := argon2.IDKey(
		[]byte(password),
		salt,
		config.Iterations,
		config.Memory,
		config.Parallelism,
		config.KeyLength,
	)

	return subtleCompare(hash, otherHash)
}

func decodeHash(encodedHash string) (*HashConfig, []byte, []byte, error) {
	var version int
	var config HashConfig

	_, err := fmt.Sscanf(
		encodedHash,
		"$argon2id$v=%d&m=%d,t=%d,p=%d",
		&version,
		&config.Memory,
		&config.Iterations,
		&config.Parallelism,
	)
	if err != nil {
		return nil, nil, nil, err
	}

	if version != argon2.Version {
		return nil, nil, nil, fmt.Errorf("incompatible argon2 version")
	}

	parts := splitHash(encodedHash)
	if len(parts) != 6 {
		return nil, nil, nil, fmt.Errorf("invalid hash format")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, err
	}
	config.SaltLength = uint32(len(salt))

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, err
	}
	config.KeyLength = uint32(len(hash))

	return &config, salt, hash, nil
}

func splitHash(encodedHash string) []string {
	var parts []string
	start := 0
	for i := 0; i < len(encodedHash); i++ {
		if encodedHash[i] == '$' {
			parts = append(parts, encodedHash[start:i])
			start = i + 1
		}
	}
	parts = append(parts, encodedHash[start:])
	return parts
}

func subtleCompare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	var result byte
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}

	return result == 0
}

func GenerateRandomString(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:length], nil
}
