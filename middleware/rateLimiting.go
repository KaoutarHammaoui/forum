package middleware

import (
	"forum/config"
	"net/http"
	"sync"
	"time"
	// "forum/handlers"
)

type requests struct {
	lastSenn time.Time
	count    int
}

var (
	visitors = make(map[string]*requests)
	mu       sync.Mutex
)

func RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		mu.Lock()
		v, exists := visitors[ip]
		if !exists {
			visitors[ip] = &requests{time.Now(), 1}
			mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		if time.Since(v.lastSenn) < time.Minute {
			if v.count >= 3 {
				mu.Unlock()
				w.WriteHeader(429)
				data:=map[string]string{
					"Error": "Too manny requests",
					"Status":http.StatusText(429),
				}
				config.RenderTemplate(w,"error.html",data)
				return 
			}
			v.count++
			} else {
			v.count = 1
			v.lastSenn = time.Now()
		}
		mu.Unlock()
		next.ServeHTTP(w, r)
	})
}
