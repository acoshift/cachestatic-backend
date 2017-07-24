package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"strings"

	"github.com/acoshift/cachestatic"
)

func main() {
	target := os.Getenv("target")
	u, err := url.Parse(target)
	if err != nil {
		panic(err)
	}
	var exclude map[string]struct{}
	{
		ps := strings.Split(os.Getenv("exclude"), "||")
		for _, p := range ps {
			exclude[p] = struct{}{}
		}
	}

	r := httputil.NewSingleHostReverseProxy(u)
	r.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Set("Server", "cachestatic-backend")
		resp.Header.Set("X-Powered-By", "acoshift")
		return nil
	}
	h := cachestatic.New(cachestatic.Config{
		Skipper: func(r *http.Request) bool {
			if _, ok := exclude[r.URL.Path]; ok {
				return true
			}
			return false
		},
	})(r)
	go func() {
		http.ListenAndServe(":8081", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "OK")
		}))
	}()
	http.ListenAndServe(":8080", h)
}
