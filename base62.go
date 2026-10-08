package main

import (
	"errors"
	"strings"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Encode turns a number into a short Base62 string
func Encode(n int64) string {
	if n == 0 {
		return "0"
	}

	var result []byte
	for n > 0 {
		result = append(result, alphabet[n%62]) // pick a character
		n = n / 62                              // move to the next "digit"
	}

	// The characters came out backwards, so flip them
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return string(result)
}

// Decode turns a Base62 string back into the original number
func Decode(s string) (int64, error) {
	var n int64
	for _, c := range []byte(s) {
		index := strings.IndexByte(alphabet, c)
		if index == -1 {
			return 0, errors.New("invalid character in code")
		}
		n = n*62 + int64(index)
	}
	return n, nil
}