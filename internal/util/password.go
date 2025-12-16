package util

import "golang.org/x/crypto/bcrypt"

// HashPassword converts a plain password into a bcrypt hash for storage.
func HashPassword(p string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// VerifyPassword compares the stored hash with the given plain password.
func VerifyPassword(hash string, p string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(p))
}
