package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/0xc00000f/shortener-tpl/internal/user"

	"github.com/0xc00000f/shortener-tpl/internal/shortener"

	"go.uber.org/zap"
)

func GetSavedData(sa *shortener.NaiveShortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := GetUserFromRequest(r)
		if !ok {
			u = user.Nil
		}
		all, err := sa.Encoder().GetAll(u.UserID)
		if err != nil {
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}

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
		w.WriteHeader(http.StatusOK)
		sa.L.Info("function result", zap.String("result", string(result)))
		w.Write(result)
	}
}

type result struct {
	Short string `json:"short_url"`
	Long  string `json:"long_url"`
}

func prepareResult(all map[string]string, baseURL string, l *zap.Logger) (b []byte, err error) {
	var res []result
	for short, long := range all {
		res = append(res, result{
			Short: fmt.Sprintf("%s/%s", baseURL, short),
			Long:  long,
		})
	}

	b, err = json.MarshalIndent(res, "", " ")
	if err != nil {
		l.Error("marshal indent error", zap.Error(err))
		return nil, err
	}
	b = append(b, '\n')
	return b, nil
}
