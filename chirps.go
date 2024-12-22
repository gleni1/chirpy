package main 

import (
  "chirpy/internal/database"
  "time"
  "net/http"
  "encoding/json"
  "context"
  "github.com/google/uuid"
  _ "github.com/lib/pq"
)

type Chirp struct {
	Body   string `json:"body"`
	UserID string `json:"user_id"`
}

func (cfg *apiConfig) handleCreateChirp(w http.ResponseWriter, r *http.Request) {
  var chirp Chirp

  const maxChirpLength = 140 
  if len(chirp.Body) > maxChirpLength {
    respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
    return 
  }
  

  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&chirp)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  cleanBodyMessage := cleanChirpBody(chirp.Body)
  user_id := uuid.MustParse(chirp.UserID)
  postParams := database.CreateChirpParams {
    ID:  uuid.New(),
    CreatedAt:  time.Now(),
    UpdatedAt:  time.Now(),
    Body:       cleanBodyMessage,
    UserID:     user_id, 
  } 
  post, err := cfg.db.CreateChirp(context.Background(), postParams)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  data, err := json.Marshal(post)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return 
  }

  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  w.WriteHeader(http.StatusCreated)
  w.Write(data)
  return
}
