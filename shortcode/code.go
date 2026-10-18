package shortcode

import (
	"crypto/rand"
	"errors"
)

// Alphabet holds the 62 characters a short code can contain.
const Alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// CodeLength is the number of characters in each short code.
const CodeLength = 7

// GenerateCode returns a random code of the given length using Alphabet.
func GenerateCode(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("code length must be greater than 0")
	}

	code := make([]byte, 0, length)
	buf := make([]byte, length+4) // a few extra, since some bytes get skipped

	for len(code) < length {
		if _, err := rand.Read(buf); err != nil {
			return "", err
		}

		for _, b := range buf {
			if b >= 248 {
				continue // skip to avoid modulo bias
			}
			code = append(code, Alphabet[b%62])
			if len(code) == length {
				break
			}
		}
	}

	return string(code), nil
}