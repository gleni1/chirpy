package main 

import (
  "fmt"
  "chirpy/internal/database"
  "chirpy/internal/auth"
  "time"
  "net/http"
  "encoding/json"
  "context"
  "github.com/google/uuid"
  _ "github.com/lib/pq"
)

type UserData struct {
  Password  string  `json:"password"`
	EmailVal  string  `json:"email"`
  Expiry    int    `json:"expires_in_seconds,omitempty"`
}

type UserResponse struct {
    ID        uuid.UUID `json:"id"`
    Email     string    `json:"email"`
    Token     string    `json:"token"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func getExpirationDuration(seconds int) time.Duration {
  if seconds <= 0 {
    return time.Hour 
  }
  if seconds > 3600 {
    return time.Hour 
  }
  return time.Duration(seconds) * time.Second
}

func (apiCfg *apiConfig) handleCreateNewUser(w http.ResponseWriter, r *http.Request) {
  var usrData UserData
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&usrData)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    fmt.Errorf("Error decoding the user data json in the request body: %w", err)
    return
  }
  emailVal := usrData.EmailVal
  unHashedPass := usrData.Password
  passwordVal, err := auth.HashPassword(unHashedPass)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    fmt.Errorf("Error hashing the password")
    return
  }

  userParams := database.CreateUserParams {
    ID:             uuid.New(),
    CreatedAt:      time.Now(),
    UpdatedAt:      time.Now(),
    Email:          emailVal, 
    HashedPassword: passwordVal,
  } 

  user, err := apiCfg.db.CreateUser(context.Background(), userParams) 
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Errorf("Error creating new user: %w", err)
    return
  }

  response := UserResponse {
    ID:           user.ID,
    Email:        user.Email,
    CreatedAt:    user.CreatedAt,
    UpdatedAt:    user.UpdatedAt,
  }

  data, err := json.Marshal(response)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Errorf("Error encoding the userParams: %w", err)
    return
  }

  w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
  w.Write(data)
  return
}

func (apiCfg *apiConfig) handleUserLogin(w http.ResponseWriter, r *http.Request) {

  var usrData UserData
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&usrData)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    fmt.Errorf("Error decoding the user data json in the request body: %w", err)
    return
  }
  emailVal := usrData.EmailVal
  unHashedPass := usrData.Password

  user, err := apiCfg.db.UserByEmail(context.Background(), emailVal)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Errorf("Error fetching user by email")
    return
  }

  err = auth.CheckPasswordHash(unHashedPass, user.HashedPassword)
  if err != nil {
    w.WriteHeader(http.StatusUnauthorized)
    return
  }

  // call the JWT function
  expirationTime := getExpirationDuration(usrData.Expiry)
  jwtToken, err := auth.MakeJWT(user.ID, apiCfg.SecretKey, expirationTime)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Errorf("Error getting token")
    return
  }
  
  response := UserResponse {
    ID:           user.ID,
    Email:        user.Email,
    Token:        jwtToken,
    CreatedAt:    user.CreatedAt,
    UpdatedAt:    user.UpdatedAt,
  }

  data, err := json.Marshal(response)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Errorf("Error encoding the userParams: %w", err)
    return
  }

  w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
  w.Write(data)
  return
}
