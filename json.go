package main 

import (
  "encoding/json"
  "log"
  "net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
  if err != nil {
    log.Println(err)
  }
  if code > 499 {
    log.Printf("Responding with 5XX error: %s", msg)
  }
  type errorResponse struct {
    Error string `json:"error"` 
  }
  respondWithJSON(w, code, errorResponse{
    Error: msg,
  })
}

// func respondWithCleanBody(w http.ResponseWriter, code int, cleanBody string) {
//   w.Header().Set("Content-Type", "application/json")
//   type cleanBodyResponse struct {
//     ClBody  string `json:"cleaned_body"`
//   }
//   respondWithJSON(w, code, cleanBodyResponse{
//     ClBody: cleanBody,
//   })
// }

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
  w.Header().Set("Content-Type", "application/json")
  dat, err := json.Marshal(payload)
  if err != nil {
    log.Printf("Error marshalling JSON: %s", err)
    w.WriteHeader(500)
    return 
  }
  w.WriteHeader(code)
  w.Write(dat)
}
