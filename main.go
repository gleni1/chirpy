package main

import (
	"log"
	"net/http"
  "sync/atomic"
  // "fmt"
)

type apiConfig struct {
  fileserverHits  atomic.Int32
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

func main() {
  var cfg = &apiConfig{}
	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.HandlerFunc(homeHandler)))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.metrics)
	mux.HandleFunc("POST /admin/reset", cfg.resetNumReq)

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
