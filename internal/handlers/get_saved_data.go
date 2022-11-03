package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"

	"github.com/0xc00000f/shortener-tpl/internal/shortener"

	"go.uber.org/zap"
)

func GetSavedData(sa *shortener.NaiveShortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		all, err := sa.Encoder().GetAll(uuid.New())
		if err != nil {
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}
		log.Printf("all: %v", all)

		if len(all) == 0 {
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		result, err := prepareResult(all, sa.BaseURL, sa.L)
		if err != nil {
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(result)
	}
}

type result struct {
	Short string `json:"short_url"`
	Long  string `json:"long_url"`
}

func prepareResult(all map[string]string, baseURL string, l *zap.Logger) ([]byte, error) {
	var res []result
	for short, long := range all {
		res = append(res, result{
			Short: fmt.Sprintf("%s/%s", baseURL, short),
			Long:  long,
		})
	}

	b, err := json.MarshalIndent(res, "", " ")
	if err != nil {
		l.Error("writing url in file marshaling error", zap.Error(err))
		return nil, err
	}
	b = append(b, '\n')
	return b, nil
}