package main

import (
	_ "database/sql"
	"net/http"
  "context"
  "time"
  "encoding/json"
  "chirpy/internal/auth"
)

func (apiCfg *apiConfig) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
  token, err := auth.GetBearerToken(r.Header)
  if err != nil {
    http.Error(w, "Unauthorized", http.StatusUnauthorized)
    return
  }

  tokenStruct, err := apiCfg.db.GetToken(context.Background(), token)
  if err != nil {
    http.Error(w, "Status Unauthorized", http.StatusUnauthorized)
    return
  }

  if time.Now().After(tokenStruct.ExpiresAt) {
    http.Error(w, "Status Unauthorized", http.StatusUnauthorized)
    return
  }

  if tokenStruct.RevokedAt.Valid {
    http.Error(w, "Status Unauthorized", http.StatusUnauthorized)
    return
  }

  accessToken, err := auth.MakeJWT(tokenStruct.UserID, apiCfg.SecretKey)
  if err != nil {
    http.Error(w, "Internal Server Error: unable to create access token", http.StatusInternalServerError)
    return
  }


  w.Header().Set("Content-Type", "application/json")
  
  response := map[string]string {
    "token": accessToken,
  }

  encoder := json.NewEncoder(w)
  err = encoder.Encode(&response)
  if err != nil {
    http.Error(w, "Internal Server Error: could not encode response", http.StatusInternalServerError)
  }

}
