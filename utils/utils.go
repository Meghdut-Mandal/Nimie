package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func ParseBody(r *http.Request, x interface{}) {
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		if err := json.Unmarshal(body, x); err != nil {
			return
		}
	}
}

//respond with error message
func RespondWithError(w http.ResponseWriter, code int, msg string) {
	RespondWithJSON(w, code, map[string]string{"error": msg})
}

func RespondWithJSONMessage(w http.ResponseWriter, code int, msg string) {
	RespondWithJSON(w, code, map[string]string{"message": msg})
}

// repond with JSON
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
