package auth

import (
  "golang.org/x/crypto/bcrypt"
  "fmt"
)

func HashPassword(password string) (string, error) {
  hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
  if err != nil {
    return "", fmt.Errorf("Error encrypting the password: %v", err)
  }
  hashedPass := string(hash)
  return hashedPass, nil 
}

func CheckPasswordHash(password, hash string) error {
  err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
  if err != nil {
    return fmt.Errorf("Passwords don't match")
  }
  return nil
}
