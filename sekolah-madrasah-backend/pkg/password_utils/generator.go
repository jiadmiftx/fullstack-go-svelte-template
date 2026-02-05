package password_utils

import (
	"crypto/rand"
	"math/big"
)

const (
	// Characters for password generation
	lowercase = "abcdefghijklmnopqrstuvwxyz"
	uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits    = "0123456789"
	// Special characters that are safe for URLs and emails
	special = "!@#$%&*"
)

// GeneratePassword creates a random password with specified length
// Default length is 12 if not specified
func GeneratePassword(length int) string {
	if length <= 0 {
		length = 12
	}

	// Combine all character sets
	allChars := lowercase + uppercase + digits + special

	// Ensure at least one character from each set
	password := make([]byte, length)

	// First 4 characters: one from each set
	if length >= 4 {
		password[0] = randomChar(lowercase)
		password[1] = randomChar(uppercase)
		password[2] = randomChar(digits)
		password[3] = randomChar(special)

		// Fill remaining with random characters from all sets
		for i := 4; i < length; i++ {
			password[i] = randomChar(allChars)
		}

		// Shuffle the password
		shuffle(password)
	} else {
		// For short passwords, just use random characters
		for i := 0; i < length; i++ {
			password[i] = randomChar(allChars)
		}
	}

	return string(password)
}

// GenerateSimplePassword creates a simpler password (letters and digits only)
func GenerateSimplePassword(length int) string {
	if length <= 0 {
		length = 8
	}

	chars := lowercase + uppercase + digits
	password := make([]byte, length)

	for i := 0; i < length; i++ {
		password[i] = randomChar(chars)
	}

	return string(password)
}

// randomChar returns a random character from the given string
func randomChar(chars string) byte {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
	return chars[n.Int64()]
}

// shuffle randomly shuffles a byte slice
func shuffle(b []byte) {
	for i := len(b) - 1; i > 0; i-- {
		j, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		b[i], b[j.Int64()] = b[j.Int64()], b[i]
	}
}
