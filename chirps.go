package main 

import (
  "chirpy/internal/database"
  "chirpy/internal/auth"
  "time"
  "log"
  "fmt"
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

  token, err := auth.GetBearerToken(r.Header)
  if err != nil {
    fmt.Printf("GetBearerToken error: %v\n", err) // Debug log
    w.WriteHeader(http.StatusUnauthorized)
    return
  }

  userID, err := auth.ValidateJWT(token, cfg.SecretKey)
  if err != nil {
    fmt.Printf("ValidateJWT error: %v\n", err) // Debug log
    w.WriteHeader(http.StatusUnauthorized)
    return
  }
  

  decoder := json.NewDecoder(r.Body)
  err = decoder.Decode(&chirp)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  cleanBodyMessage := cleanChirpBody(chirp.Body)
  //user_id := uuid.MustParse(chirp.UserID)
  postParams := database.CreateChirpParams {
    ID:  uuid.New(),
    CreatedAt:  time.Now(),
    UpdatedAt:  time.Now(),
    Body:       cleanBodyMessage,
    UserID:     userID, 
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

  authorID := r.URL.Query().Get("author_id")
  sortVal := r.URL.Query().Get("sort")

  uID, _ := uuid.Parse(authorID)

  var chirps []database.Chirp
  var err error

  if authorID != "" {
    if sortVal == "desc"{
      chirps, err = cfg.db.GetChirpsByAuthorDesc(context.Background(), uID)
      if err != nil {
        http.Error(w, "Internal server errror:", http.StatusInternalServerError)
      }
    } else {
      chirps, err = cfg.db.GetChirpsByAuthor(context.Background(), uID)
      if err != nil {
        http.Error(w, "Internal server errror:", http.StatusInternalServerError)
      }
    }
  } else {
    if sortVal == "desc" {
      chirps, err = cfg.db.GetAllChirpsDesc(context.Background())
      if err != nil {
        http.Error(w, "Internal server error:", http.StatusInternalServerError)
        return 
      }
    } else {
      chirps, err = cfg.db.GetAllChirps(context.Background())
      if err != nil {
        http.Error(w, "Internal server error:", http.StatusInternalServerError)
        return 
      }
    }
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

func (cfg *apiConfig) handleDeleteOneChirp(w http.ResponseWriter, r *http.Request) {
  chirpId := r.PathValue("chirpID")
  parsedID, err := uuid.Parse(chirpId)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest) // 400 for invalid ID format
    return
  }

  token, err := auth.GetBearerToken(r.Header)
  if err != nil {
    w.WriteHeader(http.StatusUnauthorized) // 401 for missing or invalid token
    return
  }

  userID, err := auth.ValidateJWT(token, cfg.SecretKey)
  if err != nil {
    w.WriteHeader(http.StatusUnauthorized) // 401 for failed validation
    return
  }

  // Fetch chirp details to verify ownership
  chirp, err := cfg.db.GetOneChirp(context.Background(), parsedID)
  if err != nil {
    if errors.Is(err, sql.ErrNoRows) {
      w.WriteHeader(http.StatusNotFound) // 404 if chirp not found
      return
    }
    http.Error(w, "Internal server error", http.StatusInternalServerError)
    return
  }

  if chirp.UserID != userID {
    w.WriteHeader(http.StatusForbidden) // 403 if user is not the author
    return
  }

  // Proceed with deletion if authorized
  err = cfg.db.DeleteOneChirp(context.Background(), database.DeleteOneChirpParams{
    ID: parsedID,
    UserID: userID,
  })

  if err != nil {
    // We've already checked for specific errors, remaining are internal.
    http.Error(w, "Internal server error", http.StatusInternalServerError)
    return
  }

  w.WriteHeader(http.StatusNoContent)
}
