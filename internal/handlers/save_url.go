package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/0xc00000f/shortener-tpl/internal/encoder"
	"github.com/0xc00000f/shortener-tpl/internal/shortener"
	"github.com/0xc00000f/shortener-tpl/internal/url"
	"github.com/0xc00000f/shortener-tpl/internal/user"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func SaveURL(sa *shortener.NaiveShortener) http.HandlerFunc { //revive:disable-line:cognitive-complexity
	return func(w http.ResponseWriter, r *http.Request) {
		urlPart := chi.URLParam(r, "url")

		if len(urlPart) > 0 {
			sa.L.Error("checking url param isn't success")
			http.Error(w, "400 page not found", http.StatusBadRequest)

			return
		}

		rc := r.Body
		defer rc.Close()

		u, ok := GetUserFromRequest(r)
		if !ok {
			u = user.Nil
		}

		var writeBody = func(b []byte) {
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusCreated)
			if _, err := w.Write(b); err != nil {
				sa.L.Error("writing body failure", zap.Error(err))
			}
		}

		short, err := createShort(r.Context(), sa, rc, u.UserID, false)
		if err != nil {
			var uve *encoder.UniqueViolationError

			if !errors.As(err, &uve) {
				sa.L.Error("creating short isn't success", zap.Error(err))
				http.Error(w, "400 page not found", http.StatusBadRequest)

				return
			}

			sa.L.Info("short for this long exist", zap.Error(err))

			writeBody = func(b []byte) {
				w.Header().Set("content-type", "application/json")
				w.WriteHeader(http.StatusConflict)

				if _, err := w.Write(b); err != nil {
					sa.L.Error("writing body failure", zap.Error(err))
				}
			}
		}

		body := fmt.Sprintf("%s/%s", sa.BaseURL, short)
		writeBody([]byte(body))
	}
}

type ShortRequest struct {
	URL string `json:"url"`
}

type ShortResponse struct {
	Result string `json:"result"`
}

func SaveURLJson(sa *shortener.NaiveShortener) http.HandlerFunc { //revive:disable-line:cognitive-complexity
	return func(w http.ResponseWriter, r *http.Request) {
		rc := r.Body
		defer rc.Close()

		u, ok := GetUserFromRequest(r)
		if !ok {
			u = user.Nil
		}

		var writeBody = func(b []byte) {
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusCreated)
			if _, err := w.Write(b); err != nil {
				sa.L.Error("writing body failure", zap.Error(err))
			}
		}

		short, err := createShort(r.Context(), sa, rc, u.UserID, true)
		if err != nil {
			var uve *encoder.UniqueViolationError
			if !errors.As(err, &uve) {
				sa.L.Error("creating short isn't success", zap.Error(err))
				http.Error(w, "400 page not found", http.StatusBadRequest)

				return
			}

			sa.L.Info("short for this long exist", zap.Error(err))

			writeBody = func(b []byte) {
				w.Header().Set("content-type", "application/json")
				w.WriteHeader(http.StatusConflict)

				if _, err := w.Write(b); err != nil {
					sa.L.Error("writing body failure", zap.Error(err))
				}
			}
		}

		fullEncodedURL := fmt.Sprintf("%s/%s", sa.BaseURL, short)
		resp := ShortResponse{Result: fullEncodedURL}

		respBody, err := json.Marshal(resp)
		if err != nil {
			sa.L.Error("marshalling response struct isn't success", zap.Error(err))
			http.Error(w, "400 page not found", http.StatusBadRequest)

			return
		}

		writeBody(respBody)
	}
}

func createShort(
	ctx context.Context,
	sa *shortener.NaiveShortener,
	r io.Reader,
	userID uuid.UUID,
	isJSON bool,
) (short string, err error) {
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
		return "", url.ErrInvalidURL
	}

	short, err = sa.Encoder().Short(ctx, userID, long)
	if err != nil {
		sa.L.Error("creating short isn't success", zap.Error(err))
		return short, err
	}

	return short, nil
}
