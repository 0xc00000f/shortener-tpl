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

		short, err := createShort(sa, rc, false)
		if err != nil {
			sa.L.Error("creating short isn't success", zap.Error(err))
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
		rc, err := readBody(r, sa.L)
		if err != nil {
			sa.L.Error("read body err", zap.Error(err))
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}
		defer rc.Close()

		short, err := createShort(sa, rc, true)
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

func createShort(sa *shortener.NaiveShortener, r io.Reader, isJSON bool) (short string, err error) {

	req := ShortRequest{}
	b, err := io.ReadAll(r)
	if err != nil {
		sa.L.Error("reading body isn't success", zap.Error(err))
		return "", err
	}

	var long string
	switch isJSON {
	case true:
		err = json.Unmarshal(b, &req)
		if err != nil {
			sa.L.Error("unmarshalling isn't success", zap.Error(err))
			return "", err
		}
		long = req.URL
	case false:
		long = string(b)
	}

	if !url.Valid(long) {
		sa.L.Error("url in body isn't valid")
		return "", errors.New("url in body isn't valid")
	}

	short, err = sa.Encoder().Short(long)
	if err != nil {
		sa.L.Error("creating short isn't success", zap.Error(err))
		return "", err
	}

	return short, nil
}
