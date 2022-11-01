package handlers

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/0xc00000f/shortener-tpl/internal/shortener"
	"github.com/0xc00000f/shortener-tpl/internal/url"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func SaveURL(sa *shortener.NaiveShortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlPart := chi.URLParam(r, "url")

		if len(urlPart) > 0 {
			sa.L.Error("checking url param isn't success")
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}

		rc, err := readBody(r, sa.L)
		if err != nil {
			sa.L.Error("read body err", zap.Error(err))
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}
		defer rc.Close()

		b, err := io.ReadAll(rc)
		long := string(b)
		if err != nil || !url.Valid(long) {
			sa.L.Error("checking body isn't success")
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}

		short, err := sa.Encoder().Short(long)
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

func readBody(r *http.Request, l *zap.Logger) (io.ReadCloser, error) {
	var readCloser io.ReadCloser
	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			l.Error("can't create gzip readCloser", zap.Error(err))
			return nil, errors.New("can't create gzip readCloser")
		}
		readCloser = gz
		return readCloser, nil
	}
	readCloser = r.Body

	return readCloser, nil
}

type ShortRequest struct {
	URL string `json:"url"`
}

type ShortResponse struct {
	Result string `json:"result"`
}

func SaveURLJson(sa *shortener.NaiveShortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := ShortRequest{}

		rc, err := readBody(r, sa.L)
		if err != nil {
			sa.L.Error("read body err", zap.Error(err))
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}
		defer rc.Close()

		b, err := io.ReadAll(rc)

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

		short, err := sa.Encoder().Short(req.URL)
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
