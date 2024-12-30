package main

import (
	"chirpy/internal/database"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
)

type apiConfig struct {
  fileserverHits  atomic.Int32
  db              *database.Queries
  SecretKey       string
  PolkaKey        string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    cfg.fileserverHits.Add(1)
    next.ServeHTTP(w, r)
  })
}

func (cfg *apiConfig) resetNumReq(w http.ResponseWriter, r *http.Request) {
  err := godotenv.Load()
  if err != nil {
    log.Fatalf("Error loading .env file: %w", err)
  }
  isDev := os.Getenv("PLATFORM")
  if isDev == "dev" {
    err = cfg.db.DeleteUsers(context.Background())
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      fmt.Errorf("Error deleting users: %w", err)
    }

    cfg.fileserverHits.Store(0)
    w.Write([]byte("File server hits set to 0"))
  } else {
    w.WriteHeader(http.StatusForbidden)
  }
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
  fs := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
  fs.ServeHTTP(w, r)
}


func handleResponseBody (w http.ResponseWriter, r *http.Request, msg string, stCode int) ([]byte, error) {
  var respBody map[string]interface{}

  if stCode == 200 {
    respBody = map[string]interface{} {
      "valid": true,
    }
  }
  if stCode == 400 || stCode == 500 {
    respBody = map[string]interface{} {
      "error": msg,
    }
  }

  data, err := json.Marshal(respBody)
  if err != nil {
    log.Printf("Error marshaling JSON: %s", err)
    return []byte{}, err
  }

  return data, nil
}


func main() {
  var cfg = &apiConfig{}
	const filepathRoot = "."
	const port = "8080"

  godotenv.Load()
  dbURL := os.Getenv("DB_URL")
  polka := os.Getenv("POLKA_KEY")
  //start database
  db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	defer db.Close()

  dbQueries := database.New(db)
  secKey := os.Getenv("SECRET")
  apiCfg := apiConfig{
    fileserverHits: atomic.Int32{},
    db:             dbQueries,
    SecretKey:      secKey,
    PolkaKey:       polka,
  }

	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.HandlerFunc(homeHandler)))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.metrics)
  // the below request should be a DELETE method instead
  mux.HandleFunc("POST /admin/reset", apiCfg.resetNumReq)
  // mux.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)
  mux.HandleFunc("POST /api/users", apiCfg.handleCreateNewUser)
  mux.HandleFunc("POST /api/chirps", apiCfg.handleCreateChirp)
  mux.HandleFunc("POST /api/login", apiCfg.handleUserLogin)
  mux.HandleFunc("POST /api/refresh", apiCfg.handleRefreshToken)
  mux.HandleFunc("POST /api/revoke", apiCfg.handleRevokeToken)
  mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handleWebhooks)
  mux.HandleFunc("PUT /api/users", apiCfg.handleUpdateUser)
  mux.HandleFunc("GET /api/chirps", apiCfg.handleGetChirps)
  mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handleGetOneChirp)
  mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handleDeleteOneChirp)

  srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
