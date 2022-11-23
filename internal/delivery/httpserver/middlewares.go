package httpserver

import (
	"crypto/subtle"
	"github.com/vbua/go_setter_getter/config"
	"log"
	"net/http"
	"time"
)

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%v, %v, %v, %s\n", r.Method, r.URL, timeStart.Format("2006-01-02T15:04:05"),
			time.Since(timeStart))
	})
}

func basicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user, pass, ok := r.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(config.UserName)) != 1 ||
			subtle.ConstantTimeCompare([]byte(pass), []byte(config.UserPass)) != 1 {
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
