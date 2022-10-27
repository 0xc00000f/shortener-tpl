package handlers

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"

	"github.com/0xc00000f/shortener-tpl/internal/api"
	"github.com/0xc00000f/shortener-tpl/internal/utils"
)

func SaveURL(sa api.ShortenerAPI) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlPart := chi.URLParam(r, "url")

		if len(urlPart) > 0 {
			log.Printf("checking url param isn't success")
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}

		var reader io.Reader
		if r.Header.Get(`Content-Encoding`) == `gzip` {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				log.Printf("can't create gzip reader: %v", err)
				http.Error(w, "400 page not found", http.StatusBadRequest)
				return
			}
			reader = gz
			defer gz.Close()
		} else {
			reader = r.Body
			defer r.Body.Close()
		}

		b, err := io.ReadAll(reader)
		longURL := string(b)
		if err != nil || !utils.IsURL(longURL) {
			log.Printf("checking body isn't success")
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}

		short, err := sa.Logic().Short(longURL)
		if err != nil {
			log.Printf("creating short isn't success: %v", err)
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}

		if sa.BaseURL == "" {
			sa.BaseURL = fmt.Sprintf("http://%s", r.Host)
		}

		hk := "content-type"
		hv := "text/plain; charset=utf-8"
		w.Header().Set(hk, hv)

		w.WriteHeader(http.StatusCreated)

		body := fmt.Sprintf("%s/%s", sa.BaseURL, short)
		w.Write([]byte(body))
	}
}

type ShortRequest struct {
	URL string `json:"url"`
}

type ShortResponse struct {
	Result string `json:"result"`
}

func SaveURLJson(sa api.ShortenerAPI) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := ShortRequest{}

		var reader io.Reader
		if r.Header.Get(`Content-Encoding`) == `gzip` {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				log.Printf("cant create gzip reader: %v", err)
				http.Error(w, "400 page not found", http.StatusBadRequest)
				return
			}
			reader = gz
			defer gz.Close()
		} else {
			reader = r.Body
			defer r.Body.Close()
		}

		b, err := io.ReadAll(reader)

		if err != nil {
			log.Printf("checking body isn't success: %v", err)
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}

		if err := json.Unmarshal(b, &req); err != nil {
			log.Printf("unmarshaling isn't success: %v", err)
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}

		short, err := sa.Logic().Short(req.URL)
		if err != nil {
			log.Printf("creating short isn't success: %v", err)
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}

		if sa.BaseURL == "" {
			sa.BaseURL = fmt.Sprintf("http://%s", r.Host)
		}

		fullEncodedURL := fmt.Sprintf("%s/%s", sa.BaseURL, short)
		resp := ShortResponse{Result: fullEncodedURL}

		respBody, err := json.Marshal(resp)
		if err != nil {
			log.Printf("marshalling response struct isn't success: %v", err)
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}

		hk := "content-type"
		hv := "application/json"
		w.Header().Set(hk, hv)
		w.WriteHeader(http.StatusCreated)
		w.Write(respBody)
	}
}
