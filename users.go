package main 

import (
  "fmt"
  "chirpy/internal/database"
  "time"
  "net/http"
  "encoding/json"
  "context"
  "github.com/google/uuid"
  _ "github.com/lib/pq"
)

type Email struct {
	EmailVal string `json:"email"`
}


func (apiCfg *apiConfig) handleCreateNewUser(w http.ResponseWriter, r *http.Request) {
  var email Email
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&email)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    fmt.Errorf("Error decoding the email json in the request body: %w", err)
    return
  }
  emailVal := email.EmailVal
  
  userParams := database.CreateUserParams {
    ID:           uuid.New(),
    CreatedAt:    time.Now(),
    UpdatedAt:    time.Now(),
    Email:        emailVal, 
  } 


  user, err := apiCfg.db.CreateUser(context.Background(), userParams) 
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Errorf("Error creating new user: %w", err)
    return
  }

  data, err := json.Marshal(user)
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


