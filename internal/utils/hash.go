package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
)

const (
	// Argon2id configuration
	argonTime    uint32 = 1
	argonMemory  uint32 = 64 * 1024 // 64 MB
	argonThreads uint8  = 4
	argonKeyLen  uint32 = 32

	saltLength = 16
)

// Hash format:
//
//	argon2id$time$memory$threads$salt$hash
//
// Example:
//
//	argon2id$1$65536$4$base64salt$base64hash

func HashPassword(password string) (string, error) {

	salt := make([]byte, saltLength)

	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf(
			"failed to generate salt: %w",
			err,
		)
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		argonTime,
		argonMemory,
		argonThreads,
		argonKeyLen,
	)

	saltEncoded := base64.RawStdEncoding.EncodeToString(salt)

	hashEncoded := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf(
		"argon2id$%d$%d$%d$%s$%s",
		argonTime,
		argonMemory,
		argonThreads,
		saltEncoded,
		hashEncoded,
	)

	return encodedHash, nil
}

// CheckPassword verifies password.
//
// Returns:
//
//	valid          -> password matched
//	rehashRequired -> bcrypt user should migrate to argon2id
//	newHash        -> newly generated argon2id hash
func CheckPassword(
	password string,
	storedHash string,
) (
	valid bool,
	rehashRequired bool,
	newHash string,
	err error,
) {

	// bcrypt backward compatibility
	if isBcryptHash(storedHash) {

		valid = checkBcryptPassword(
			password,
			storedHash,
		)

		if !valid {
			return false, false, "", nil
		}

		// migrate bcrypt -> argon2id
		rehash, err := HashPassword(password)

		if err != nil {
			return true, false, "", fmt.Errorf(
				"failed to rehash bcrypt password: %w",
				err,
			)
		}

		return true, true, rehash, nil
	}

	// argon2id verification
	valid = checkArgon2Password(
		password,
		storedHash,
	)

	return valid, false, "", nil
}

func NeedsRehash(storedHash string) bool {

	// bcrypt hashes should always migrate
	if isBcryptHash(storedHash) {
		return true
	}

	parts := strings.Split(storedHash, "$")

	if len(parts) != 6 {
		return true
	}

	timeCost, err := strconv.ParseUint(
		parts[1],
		10,
		32,
	)

	if err != nil {
		return true
	}

	memoryCost, err := strconv.ParseUint(
		parts[2],
		10,
		32,
	)

	if err != nil {
		return true
	}

	threadCount, err := strconv.ParseUint(
		parts[3],
		10,
		8,
	)

	if err != nil {
		return true
	}

	return uint32(timeCost) != argonTime ||
		uint32(memoryCost) != argonMemory ||
		uint8(threadCount) != argonThreads
}

func checkArgon2Password(
	password string,
	storedHash string,
) bool {

	params, salt, hash, err := parseArgon2Hash(storedHash)

	if err != nil {
		return false
	}

	computedHash := argon2.IDKey(
		[]byte(password),
		salt,
		params.Time,
		params.Memory,
		params.Threads,
		uint32(len(hash)),
	)

	return subtle.ConstantTimeCompare(
		computedHash,
		hash,
	) == 1
}

func checkBcryptPassword(
	password string,
	hash string,
) bool {

	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	) == nil
}

func isBcryptHash(hash string) bool {

	return strings.HasPrefix(hash, "$2a$") ||
		strings.HasPrefix(hash, "$2b$") ||
		strings.HasPrefix(hash, "$2y$")
}

type argon2Params struct {
	Time    uint32
	Memory  uint32
	Threads uint8
}

func parseArgon2Hash(
	encodedHash string,
) (*argon2Params, []byte, []byte, error) {

	parts := strings.Split(encodedHash, "$")

	if len(parts) != 6 {
		return nil, nil, nil,
			errors.New("invalid hash format")
	}

	if parts[0] != "argon2id" {
		return nil, nil, nil,
			errors.New("unsupported hash type")
	}

	timeCost, err := strconv.ParseUint(
		parts[1],
		10,
		32,
	)

	if err != nil {
		return nil, nil, nil,
			fmt.Errorf("invalid time cost: %w", err)
	}

	memoryCost, err := strconv.ParseUint(
		parts[2],
		10,
		32,
	)

	if err != nil {
		return nil, nil, nil,
			fmt.Errorf("invalid memory cost: %w", err)
	}

	threadCount, err := strconv.ParseUint(
		parts[3],
		10,
		8,
	)

	if err != nil {
		return nil, nil, nil,
			fmt.Errorf("invalid thread count: %w", err)
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])

	if err != nil {
		return nil, nil, nil,
			fmt.Errorf("invalid salt: %w", err)
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])

	if err != nil {
		return nil, nil, nil,
			fmt.Errorf("invalid hash: %w", err)
	}

	params := &argon2Params{
		Time:    uint32(timeCost),
		Memory:  uint32(memoryCost),
		Threads: uint8(threadCount),
	}

	return params, salt, hash, nil
}
