package middleware

import (
	"net/http"
	"sync"
	"time"
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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		ip:=r.RemoteAddr
		mu.Lock()
		v,exists:=visitors[ip]
		if !exists{
			visitors[ip]=&requests{time.Now(),1}
			mu.Unlock()
			next.ServeHTTP(w,r)
			return 
		}
		if time.Since(v.lastSenn) <time.Minute{
			if v.count >=20{
				mu.Unlock()
				http.Error(w,"too many request",429)
				return
			}
			v.count++
		}else{
			v.count=1
			v.lastSenn=time.Now()
		}
		mu.Unlock()
		next.ServeHTTP(w,r)
	})
}
