package main 

import (
  // "encoding/json"
  // "net/http"
  "strings"
)

type chirp struct {
  Body  string  `json:"body"`
}

type returnVal struct {
  BodyMsg   string   `json:"cleaned_body"`
}


func cleanChirpBody(msgBody string) string {
  badWords := map[string]struct{}{
    "kerfuffle": {},
    "sharbert": {},
    "fornax": {},
  }
  words := strings.Fields(msgBody)

  for i, word := range words {
    if _, exists := badWords[strings.ToLower(word)]; exists {
      words[i] = "****"
    } 
  }
  cleanMsg := strings.Join(words, " ")
  return cleanMsg
}
