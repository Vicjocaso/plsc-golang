package lib

import (
	"encoding/json"
	"net/http"
)

func EncodeOkReponse(w http.ResponseWriter, i interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(i)
}
