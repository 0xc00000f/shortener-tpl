package handlers

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/0xc00000f/shortener-tpl/internal/api"
	"github.com/0xc00000f/shortener-tpl/internal/url"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func SaveURL(sa *api.ShortenerAPI) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlPart := chi.URLParam(r, "url")

		if len(urlPart) > 0 {
			sa.L.Error("checking url param isn't success")
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}

		var reader io.Reader
		if r.Header.Get(`Content-Encoding`) == `gzip` {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				sa.L.Error("can't create gzip reader", zap.Error(err))
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
		long := string(b)
		if err != nil || !url.Valid(long) {
			sa.L.Error("checking body isn't success")
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}

		short, err := sa.Logic().Short(long)
		if err != nil {
			sa.L.Error("creating short isn't success: %v", zap.Error(err))
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

func SaveURLJson(sa *api.ShortenerAPI) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := ShortRequest{}

		var reader io.Reader
		if r.Header.Get(`Content-Encoding`) == `gzip` {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				sa.L.Error("cant create gzip reader", zap.Error(err))
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
			sa.L.Error("checking body isn't success", zap.Error(err))
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}

		if err := json.Unmarshal(b, &req); err != nil {
			sa.L.Error("unmarshalling isn't success", zap.Error(err))
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}

		short, err := sa.Logic().Short(req.URL)
		if err != nil {
			sa.L.Error("creating short isn't success", zap.Error(err))
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
			sa.L.Error("marshalling response struct isn't success", zap.Error(err))
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
