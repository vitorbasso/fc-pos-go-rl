package ratel

import (
	"log"
	"net"
	"net/http"
)

const (
	overLimitMessage = "you have reached the maximum number of requests or actions allowed within a certain time frame"
)

func Middleware(options ...Option) func(next http.Handler) http.Handler {
	ratel := NewLimiter(options...)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			k := r.Header.Get("API_KEY")
			if k == "" {
				k = getIP(r)
			}
			allowed, err := ratel.Allow(k)
			if err != nil {
				log.Printf("error checking rate limit: %v", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			if !allowed {
				http.Error(w, overLimitMessage, http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func getIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	var err error
	if ip == "" {
		ip, _, err = net.SplitHostPort(r.RemoteAddr)
	}
	if err != nil || ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}
