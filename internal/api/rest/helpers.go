package rest

import (
	"encoding/json"
	"net/http"
)

const (
	ApplicationJSONContentType = "application/json; charset=utf-8"
)

func JSONError(w http.ResponseWriter, message string, code int) {
	response := map[string]string{"error": message}
	w.Header().Set("Content-Type", ApplicationJSONContentType)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(response)
}
