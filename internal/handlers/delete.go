package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/0xc00000f/shortener-tpl/internal/models"
	"github.com/0xc00000f/shortener-tpl/internal/shortener"
	"github.com/0xc00000f/shortener-tpl/internal/user"
)

func Delete(sa *shortener.NaiveShortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("handling delete")
		now := time.Now()

		u, ok := GetUserFromRequest(r)
		if !ok {
			u = user.Nil
		}

		rc := r.Body
		defer rc.Close()

		b, err := io.ReadAll(rc)
		if err != nil {
			sa.L.Error("reading body isn't success", zap.Error(err))
			http.Error(w, "400 page not found", http.StatusBadRequest)

			return
		}

		ib, err := parseDeleteResp(b)
		if err != nil {
			sa.L.Error("unmarshalling isn't success", zap.Error(err))
			http.Error(w, "400 page not found", http.StatusBadRequest)

			return
		}

		chunkSize := 10
		chunks := chunkSlice(ib.Array, chunkSize)

		log.Printf("input data: %s", ib.Array)

		go func() {
			for i := 0; i < len(chunks); i++ {
				currentChunk := chunks[i]

				go func() {
					sa.Job <- DeleteJob{sa: sa, urlChunk: short2url(u.UserID, currentChunk)}
				}()
			}
		}()

		log.Printf("delete handled: %v", time.Since(now))

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusAccepted)
	}
}

type DeleteResp struct {
	Array []string
}

func parseDeleteResp(b []byte) (dr DeleteResp, err error) {
	err = json.Unmarshal(b, &dr.Array)
	return dr, err
}

type DeleteJob struct {
	sa       *shortener.NaiveShortener
	urlChunk []models.URL
}

func (j DeleteJob) Run(ctx context.Context) error {
	return j.sa.Encoder().Delete(ctx, j.urlChunk)
}

func chunkSlice(slice []string, chunkSize int) [][]string {
	var chunks [][]string
	for {
		if len(slice) == 0 {
			break
		}

		// necessary check to avoid slicing beyond
		// slice capacity
		if len(slice) < chunkSize {
			chunkSize = len(slice)
		}

		chunks = append(chunks, slice[0:chunkSize])
		slice = slice[chunkSize:]
	}

	return chunks
}

func short2url(uid uuid.UUID, shorts []string) []models.URL {
	urls := make([]models.URL, 0, len(shorts))

	for _, short := range shorts {
		urls = append(urls, models.URL{
			UserID: uid,
			Short:  short,
			Long:   "",
		})
	}

	return urls
}
