package main 

import (
	"net/http"
  "html/template"
  // "fmt"
)

type Data struct {
  Hits  int32
}

func (cfg *apiConfig) metrics(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  
  tmpl, err := template.ParseFiles("./metrics.html")
  if err != nil {
    http.Error(w, "Error parsing template", http.StatusInternalServerError)
    return 
  }

  var data = &Data{}
  data.Hits = cfg.fileserverHits.Load()

  if err := tmpl.Execute(w, data); err != nil {
    http.Error(w, "Error rendering template", http.StatusInternalServerError)
  }
}
