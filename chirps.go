package main 

import (
  "chirpy/internal/database"
  "time"
  "log"
  "errors"
  "database/sql"
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
	
func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    return
  }

  chirps, err := cfg.db.GetAllChirps(context.Background())
  if err != nil {
    http.Error(w, "Internal server error:", http.StatusInternalServerError)
    return 
  }


  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  w.WriteHeader(http.StatusOK)

  encoder := json.NewEncoder(w)
  err = encoder.Encode(chirps)
  if err != nil {
    log.Printf("Error encoding response: %w", err)
    return
  }
}


func (cfg *apiConfig) handleGetOneChirp(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    return
  }

  chirpId := r.PathValue("chirpID")
  parsedID, err := uuid.Parse(chirpId)
  if err != nil {
    w.WriteHeader(http.StatusNotFound)
    return
  }
  chirp, err := cfg.db.GetOneChirp(context.Background(), parsedID)
  if err != nil {
    if errors.Is(err, sql.ErrNoRows) {
      w.WriteHeader(http.StatusNotFound)
      return
    }
    http.Error(w, "Internal server error", http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  w.WriteHeader(http.StatusOK)

  encoder := json.NewEncoder(w)
  err = encoder.Encode(chirp)
  if err != nil {
    log.Printf("Error encoding response: %v", err)
    return
  }
}
