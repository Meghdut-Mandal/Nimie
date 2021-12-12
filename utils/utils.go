package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
)

func ParseBody(r *http.Request, x interface{}) {
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		if err := json.Unmarshal(body, x); err != nil {
			return
		}
	}
}

// ParseInt64 Convert string to int64
func ParseInt64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		fmt.Errorf("error: %s", err)
		return 0
	}
	return i
}

// GetUserId get user id from request
func GetUserId(r *http.Request) int64 {
	userId := int64(0)
	r.Header.Get("user_id")
	if userIdStr := r.Header.Get("user_id"); userIdStr != "" {
		userId = ParseInt64(userIdStr)
	}
	return userId
}

var (
	// CharacterSet consists of 62 characters [0-9][A-Z][a-z].
	CharacterSet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	UnixTime2021 = int64(16094592001000000)
)

// Base62Encode Encode returns a base62 representation as  string of the given integer number.
func Base62Encode(num int64) string {
	b := make([]byte, 0)

	// loop as long the num is bigger than zero
	for num > 0 {
		// receive the rest
		r := math.Mod(float64(num), float64(62))

		// devide by Base
		num /= 62

		// append chars
		b = append([]byte{CharacterSet[int(r)]}, b...)
	}

	return string(b)
}

// RespondWithError respond with error message
func RespondWithError(w http.ResponseWriter, code int, msg string) {
	RespondWithJSON(w, code, map[string]string{"error": msg})
}

func RespondWithJSONMessage(w http.ResponseWriter, code int, msg string) {
	RespondWithJSON(w, code, map[string]string{"message": msg})
}

// RespondWithJSON repond with JSON
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
