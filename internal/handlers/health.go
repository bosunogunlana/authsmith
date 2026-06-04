package handlers

import (
	"fmt"
	"net/http"
)

func (h *Handlers) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	fmt.Fprintf(w, `{"status":"ok"}`)
}