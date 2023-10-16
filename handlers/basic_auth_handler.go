package handlers

import (
	"crypto/subtle"
	"net/http"
	"short-link/config"
)

func BasicAuth(w http.ResponseWriter, r *http.Request) bool {

	const realm = "Provide user name and password"

	user, pass, ok := r.BasicAuth()

	if !ok ||
		subtle.ConstantTimeCompare([]byte(user), []byte(config.AppConfig.Username)) != 1 ||
		subtle.ConstantTimeCompare([]byte(pass), []byte(config.AppConfig.Password)) != 1 {

		w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
		w.WriteHeader(401)

		_, _ = w.Write([]byte("Unauthorised.\n"))
		return false
	}

	return true
}

func Authorize(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if BasicAuth(w, r) {
			handler(w, r)
		}
	}
}
