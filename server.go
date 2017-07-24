package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/acoshift/cachestatic"
)

func main() {
	target := os.Getenv("target")
	u, err := url.Parse(target)
	if err != nil {
		panic(err)
	}
	r := httputil.NewSingleHostReverseProxy(u)
	r.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Set("Server", "cachestatic-backend")
		resp.Header.Set("X-Powered-By", "acoshift")
		return nil
	}
	h := cachestatic.New(cachestatic.DefaultConfig)(r)
	go func() {
		http.ListenAndServe(":8081", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "OK")
		}))
	}()
	http.ListenAndServe(":8080", h)
}
