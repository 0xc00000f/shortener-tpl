package handlers

import (
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
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}

		b, err := io.ReadAll(r.Body)
		longURL := string(b)
		if err != nil || !utils.IsURL(longURL) {
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}

		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		short, err := sa.Logic().Short(longURL)
		if err != nil {
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}

		if sa.BaseURL == "" {
			sa.BaseURL = fmt.Sprintf("http://%s", r.Host)
		}

		//
		//baseURL, ok := os.LookupEnv("BASE_URL")
		//var fullEncodedURL string
		//if ok {
		//	fullEncodedURL = fmt.Sprintf("%s/%s", baseURL, short)
		//} else {
		//	fullEncodedURL = fmt.Sprintf("http://%s/%s", r.Host, short)
		//}
		w.Write([]byte(fmt.Sprintf("%s/%s", sa.BaseURL, short)))
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
		b, err := io.ReadAll(r.Body)
		log.Printf("b:%v", b)
		log.Printf("b string:%v", string(b))
		if err != nil {
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if err := json.Unmarshal(b, &req); err != nil {
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}
		log.Printf("req:%v", req)

		short, err := sa.Logic().Short(req.URL)
		if err != nil {
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}
		log.Printf("short:%v", short)

		if sa.BaseURL == "" {
			sa.BaseURL = fmt.Sprintf("http://%s", r.Host)
		}
		fullEncodedURL := fmt.Sprintf("%s/%s", sa.BaseURL, short)
		resp := ShortResponse{Result: fullEncodedURL}

		respBody, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}
		log.Printf("respBody:%v", respBody)
		log.Printf("respBodyString:%v", string(respBody))

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(respBody)
	}
}
