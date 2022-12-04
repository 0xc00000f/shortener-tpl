package handlers

import (
	"compress/gzip"
	"net/http"
)

func UnzipBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Header.Get(`Content-Encoding`) == `gzip` {
			gz, err := gzip.NewReader(req.Body)
			if err == nil {
				req.Body = gz
			}
		}

		next.ServeHTTP(w, req)
	})
}
