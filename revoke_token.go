package main

import (
	_ "database/sql"
	"net/http"
  "context"
  "chirpy/internal/auth"
)

func (apiCfg *apiConfig) handleRevokeToken(w http.ResponseWriter, r *http.Request) {
  token, err := auth.GetBearerToken(r.Header)
  if err != nil {
    http.Error(w, "Unauthorized", http.StatusUnauthorized)
    return
  }

  err = apiCfg.db.RevokeToken(context.Background(), token)
  if err != nil {
    http.Error(w, "Unauthorized", http.StatusUnauthorized)
    return
  }

  w.WriteHeader(http.StatusNoContent)
}
