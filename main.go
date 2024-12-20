package main

import (
	"log"
	"net/http"
  "sync/atomic"
  "encoding/json"
  // "fmt"
)

type apiConfig struct {
  fileserverHits  atomic.Int32
}

type returnVal struct {
  Msg   string    `json:"msg"`
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    cfg.fileserverHits.Add(1)
    next.ServeHTTP(w, r)
  })
}

func (cfg *apiConfig) resetNumReq(w http.ResponseWriter, r *http.Request) {
  cfg.fileserverHits.Store(0)
  w.Write([]byte("File server hits set to 0"))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
  fs := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
  fs.ServeHTTP(w, r)
}

func handleIncomingJson(w http.ResponseWriter, r *http.Request) {
  type chirp struct {
    Body  string  `json:"body"`
  }

  newChirp := chirp{}
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&newChirp)
  if err != nil {
    // log.Printf("Error decoding the chirp: %s", err)
    msg := "Something went wrong"
    data, _ := handleResponseBody(w, r, msg, 500)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(500)
    w.Write(data)
    return
  }

  lenChirp := len(newChirp.Body)
  if lenChirp > 140 {
    msg := "Chirp is too long"
    data, _ := handleResponseBody(w, r, msg, 400)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(400)
    w.Write(data)
    return 
  }
  msg := "no_purpose"
  data, _ := handleResponseBody(w, r, msg, 200)
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(200)
  w.Write(data)
  return
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

	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.HandlerFunc(homeHandler)))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.metrics)
	mux.HandleFunc("POST /admin/reset", cfg.resetNumReq)
  mux.HandleFunc("POST /api/validate_chirp", handleIncomingJson)

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
