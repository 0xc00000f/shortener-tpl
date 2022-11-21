package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/0xc00000f/shortener-tpl/internal/shortener"
	"github.com/0xc00000f/shortener-tpl/internal/user"
	"go.uber.org/zap"
)

func Batch(sa *shortener.NaiveShortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := GetUserFromRequest(r)
		if !ok {
			u = user.Nil
		}

		rc, err := unzipBody(r, sa.L)
		if err != nil {
			sa.L.Error("read body err", zap.Error(err))
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}
		defer rc.Close()

		b, err := io.ReadAll(rc)
		if err != nil {
			sa.L.Error("reading body isn't success", zap.Error(err))
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}

		ib, err := parseInputBatch(b)
		if err != nil {
			sa.L.Error("unmarshalling isn't success", zap.Error(err))
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}

		result, err := prepareOutputBatchResult(ib, sa, u)
		if err != nil {
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		sa.L.Info("function result", zap.String("result", string(result)))
		w.Write(result)
	}
}

type inputBatch struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type outputBatch struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func parseInputBatch(b []byte) (ib []inputBatch, err error) {
	err = json.Unmarshal(b, &ib)
	return
}

func prepareOutputBatchResult(ib []inputBatch, sa *shortener.NaiveShortener, u user.User) (result []byte, err error) {
	var ob []outputBatch
	for _, batch := range ib {
		short, err := sa.Encoder().Short(u.UserID, batch.OriginalURL)
		if err != nil {
			sa.L.Error("batch creating short isn't success", zap.Error(err))
			return nil, err
		}

		ob = append(ob, outputBatch{
			CorrelationID: batch.CorrelationID,
			ShortURL:      fmt.Sprintf("%s/%s", sa.BaseURL, short),
		})
	}

	result, err = json.MarshalIndent(ob, "", " ")
	if err != nil {
		sa.L.Error("batch marshalling isn't success", zap.Error(err))
	}
	return
}
