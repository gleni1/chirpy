package auth

import (
  "golang.org/x/crypto/bcrypt"
  "fmt"
  "time"
  "errors"
  "strings"
  "crypto/rand"
  "encoding/hex"
  "net/http"
  "github.com/google/uuid"
  "github.com/golang-jwt/jwt/v5"
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
// MakeJWT generates a signed JWT for the given userID with a specified expiration time.
func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {
	// Validate userID
	if userID == uuid.Nil {
		return "", errors.New("invalid UUID: userID cannot be empty")
	}

	// Define claims
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
    ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		Subject:   userID.String(),
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the provided secret
	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("error signing the token: %v", err)
	}

	return signedToken, nil
}

// ValidateJWT parses and validates the provided JWT, returning the userID from the token.
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}

	// Parse token with claims
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(tokenSecret), nil
		},
	)

	// Return error if parsing fails
	if err != nil {
		return uuid.Nil, fmt.Errorf("error parsing token: %v", err)
	}

	// Check if token is valid
	if !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	// Validate the subject (user ID) field
	if claims.Subject == "" {
		return uuid.Nil, errors.New("token subject is missing")
	}

	// Parse and return the user ID
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID in token subject: %v", err)
	}

	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
  authHeader := headers.Get("Authorization")

  if authHeader == "" {
    return "", errors.New("authorization header not present")
  }

  const bearerPrefix = "Bearer "
  if !strings.HasPrefix(authHeader, bearerPrefix) {
    return "", errors.New("malformed authorization header")
  }

  token := strings.TrimPrefix(authHeader, bearerPrefix)

  return strings.TrimSpace(token), nil
}

func MakeRefreshToken() (string, error) {
  randString := make([]byte, 32)
  _, err := rand.Read(randString)
  if err != nil {
    return "", err
  }
  finalString := hex.EncodeToString(randString)
  return finalString, nil
}


func GetAPIKey(headers http.Header) (string, error) {
  authHeader := headers.Get("Authorization")

  if authHeader == "" {
    return "", errors.New("authorization header not present")
  }

  const apiPrefix = "ApiKey "
  if !strings.HasPrefix(authHeader, apiPrefix) {
    return "", errors.New("malformed authorization header")
  }

  apiKey := strings.TrimPrefix(authHeader, apiPrefix)

  return strings.TrimSpace(apiKey), nil
}
