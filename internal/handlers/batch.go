package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/0xc00000f/shortener-tpl/internal/shortener"
	"github.com/0xc00000f/shortener-tpl/internal/user"
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

		result, err := prepareOutputBatchResult(r.Context(), ib, sa, u)
		if err != nil {
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		sa.L.Info("function result", zap.String("result", string(result)))

		if _, err := w.Write(result); err != nil {
			sa.L.Error("writing body failure", zap.Error(err))
		}
	}
}

type inputBatch struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type OutputBatch struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func parseInputBatch(b []byte) (ib []inputBatch, err error) {
	err = json.Unmarshal(b, &ib)
	return ib, err
}

func prepareOutputBatchResult(
	ctx context.Context,
	ib []inputBatch,
	sa *shortener.NaiveShortener,
	u user.User,
) (result []byte, err error) {
	ob := make([]OutputBatch, 0, len(ib))

	for _, batch := range ib {
		short, err := sa.Encoder().Short(ctx, u.UserID, batch.OriginalURL)
		if err != nil {
			sa.L.Error("batch creating short isn't success", zap.Error(err))
			return nil, err
		}

		ob = append(ob, OutputBatch{
			CorrelationID: batch.CorrelationID,
			ShortURL:      fmt.Sprintf("%s/%s", sa.BaseURL, short),
		})
	}

	result, err = json.MarshalIndent(ob, "", " ")
	if err != nil {
		sa.L.Error("batch marshalling isn't success", zap.Error(err))
	}

	return result, err
}
