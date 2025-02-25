package http

import (
	"encoding/json"
	"mesaggio-test/internal/models"
	"net/http"
)

func Respond(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func RespondError(w http.ResponseWriter, code int, err error) {
	Respond(w, code, models.ErrorResponse{Error: err.Error()})
}
