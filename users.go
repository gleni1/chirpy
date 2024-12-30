package main 

import (
  "fmt"
  "chirpy/internal/database"
  "chirpy/internal/auth"
  "time"
  "net/http"
  "encoding/json"
  "database/sql"
  "context"
  "github.com/google/uuid"
  _ "github.com/lib/pq"
)

type UserData struct {
  Password  string  `json:"password"`
	EmailVal  string  `json:"email"`
  // Expiry    int    `json:"expires_in_seconds,omitempty"`
}

type UserResponse struct {
    ID            uuid.UUID `json:"id"`
    Email         string    `json:"email"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
    Token         string    `json:"token"`
    RefreshToken  string    `json:"refresh_token"`
    IsChirpyRed   bool      `json:"is_chirpy_red"`
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
    IsChirpyRed:  user.IsChirpyRed,
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

  jwtToken, err := auth.MakeJWT(user.ID, apiCfg.SecretKey)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Errorf("Error getting token")
    return
  }
 
  refreshTokenString, err := auth.MakeRefreshToken()
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Errorf("Error getting refresh token string")
    return
  }

  // calculate expirationTime for refresh_token
  expiryDate := time.Now().Add(60 * 24 * time.Hour) // 60 days from now

  // call the function to create new entry in the table for the refresh_token
  refreshToken, err := apiCfg.db.CreateRefreshToken(context.Background(), database.CreateRefreshTokenParams{
    Token:      refreshTokenString,
    CreatedAt:  time.Now(),
    UpdatedAt:  time.Now(),
    UserID:     user.ID,
    ExpiresAt:  expiryDate,
    RevokedAt:  sql.NullTime{},
  })
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Errorf("Error creating refresh token")
    return
  }

  response := UserResponse {
    ID:           user.ID,
    Email:        user.Email,
    Token:        jwtToken,
    RefreshToken: refreshToken.Token,
    CreatedAt:    user.CreatedAt,
    UpdatedAt:    user.UpdatedAt,
    IsChirpyRed:  user.IsChirpyRed,
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



func (apiCfg *apiConfig) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
  token, err := auth.GetBearerToken(r.Header)
  if err != nil {
    w.WriteHeader(http.StatusUnauthorized)
    return
  }
  
  userID, err := auth.ValidateJWT(token, apiCfg.SecretKey)
  if err != nil {
    w.WriteHeader(http.StatusUnauthorized)
    return
  }


  var usrData UserData
  decoder := json.NewDecoder(r.Body)
  err = decoder.Decode(&usrData)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Errorf("Error decoding the userParams: %w", err)
    return
  }

  password := usrData.Password
  hashedPassword, err := auth.HashPassword(password)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Errorf("Error hashing the password")
    return
  }


  err = apiCfg.db.UpdateUser(context.Background(), database.UpdateUserParams{
    Email:              usrData.EmailVal,
    HashedPassword:     hashedPassword,
    ID:                 userID,
  })

  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Errorf("Error updating the user info")
    return
  }

  updatedUser := struct {
    Email   string    `json:"email"` 
  }{
    Email:  usrData.EmailVal,
  }


  w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(updatedUser)
  return
  
}

type Webhook struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

func (apiCfg *apiConfig) handleWebhooks(w http.ResponseWriter, r *http.Request) {
  var webhook Webhook
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&webhook)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  if webhook.Event != "user.upgraded" {
    w.WriteHeader(http.StatusNoContent)
    return
  }

  userID, _ := uuid.Parse(webhook.Data.UserID)

  _, err = apiCfg.db.GetUserById(context.Background(), userID)
  if err != nil {
    w.WriteHeader(http.StatusNotFound)
    return
  }

  err = apiCfg.db.UpgradeToChirpyRed(context.Background(), userID)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  w.WriteHeader(http.StatusNoContent)
  return

}
