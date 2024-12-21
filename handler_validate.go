package main 

import (
  "encoding/json"
  "net/http"
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


func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
  var chrp chirp
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&chrp)
  
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
    return 
  }

  const maxChirpLength = 140 
  if len(chrp.Body) > maxChirpLength {
    respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
    return 
  }
  
  cleanedBody := cleanChirpBody(chrp.Body)
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, "Couldn't clean chirp body", err)
    return
  }

  respondWithJSON(w, http.StatusOK, returnVal {
    BodyMsg: cleanedBody,
  })
}
