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
		log.Printf("handlers::SaveURL -- entered, arguments - sa: %v", sa)
		defer log.Print("handlers::SaveURL -- finished")

		log.Print("handlers::SaveURL -- checking correct url param")
		urlPart := chi.URLParam(r, "url")
		log.Printf("handlers::SaveURL -- checking correct url param, urlPart: %s", urlPart)

		if len(urlPart) > 0 {
			log.Printf("handlers::SaveURL -- checking url param isn't success, returning 400 err")
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}

		b, err := io.ReadAll(r.Body)
		log.Printf("handlers::SaveURL -- reading body, body: %v, err: %v", string(b), err)
		longURL := string(b)
		if err != nil || !utils.IsURL(longURL) {
			log.Printf("handlers::SaveURL -- checking is body url or err != nil isn't success, returning 400 err")
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		log.Printf("handlers::SaveURL -- call Logic::Short, params -- longURL: %v", longURL)
		short, err := sa.Logic().Short(longURL)
		log.Printf("handlers::SaveURL -- call Logic::Short, returned -- short: %s, err: %v", short, err)
		if err != nil {
			log.Printf("handlers::SaveURL -- call Logic::Short isn't success, returning 400 err")
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}

		log.Print("handlers::SaveURL -- checking sa.BaseURL is inited")
		if sa.BaseURL == "" {
			sa.BaseURL = fmt.Sprintf("http://%s", r.Host)
			log.Print("handlers::SaveURL -- sa.BaseURL isn't inited")
			log.Printf("handlers::SaveURL -- sa.BaseURL init with host value -- sa.BaseURL: %s", sa.BaseURL)
		}

		log.Print("handlers::SaveURL -- write in ResponseWriter")

		hk := "content-type"
		hv := "text/plain; charset=utf-8"
		w.Header().Set(hk, hv)
		log.Printf("handlers::SaveURL -- write in ResponseWriter, header key: %s, header value: %s", hk, hv)

		w.WriteHeader(http.StatusCreated)
		log.Printf("handlers::SaveURL -- write in ResponseWriter, status: %v", http.StatusCreated)

		body := fmt.Sprintf("%s/%s", sa.BaseURL, short)
		w.Write([]byte(body))
		log.Printf("handlers::SaveURL -- write in ResponseWriter, body: %s", body)
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
		log.Printf("handlers::SaveURLJson -- entered, arguments - sa: %v", sa)
		defer log.Print("handlers::SaveURLJson -- finished")

		req := ShortRequest{}
		b, err := io.ReadAll(r.Body)
		log.Printf("handlers::SaveURLJson -- reading body, body: %v, err: %v", string(b), err)

		if err != nil {
			log.Printf("handlers::SaveURLJson -- checking is body url or err != nil isn't success, returning 400 err")
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		log.Printf("handlers::SaveURLJson -- unmarshaling body")
		if err := json.Unmarshal(b, &req); err != nil {
			log.Printf("handlers::SaveURLJson -- unmarshaling isn't success, returning 400 err")
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}
		log.Printf("handlers::SaveURLJson -- unmarshaling is success, req: %v", req)

		log.Printf("handlers::SaveURLJson -- call Logic::Short, params -- longURL: %v", req.URL)
		short, err := sa.Logic().Short(req.URL)
		log.Printf("handlers::SaveURLJson -- call Logic::Short, returned -- short: %s, err: %v", short, err)
		if err != nil {
			log.Printf("handlers::SaveURLJson -- call Logic::Short isn't success, returning 400 err")
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}

		log.Print("handlers::SaveURLJson -- checking sa.BaseURL is inited")
		if sa.BaseURL == "" {
			sa.BaseURL = fmt.Sprintf("http://%s", r.Host)
			log.Print("handlers::SaveURLJson -- sa.BaseURL isn't inited")
			log.Printf("handlers::SaveURLJson -- sa.BaseURL init with host value -- sa.BaseURL: %s", sa.BaseURL)
		}

		log.Printf("handlers::SaveURLJson -- creating response struct")
		fullEncodedURL := fmt.Sprintf("%s/%s", sa.BaseURL, short)
		resp := ShortResponse{Result: fullEncodedURL}
		log.Printf("handlers::SaveURLJson -- creating response struct, resp: %v", resp)

		log.Printf("handlers::SaveURLJson -- marshalling response struct")
		respBody, err := json.Marshal(resp)
		log.Printf("handlers::SaveURLJson -- marshalling response struct, "+
			"respBody: %v, err: %v", string(respBody), err)
		if err != nil {
			log.Printf("handlers::SaveURLJson -- marshalling response struct isn't success, returning 400 err")
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}

		log.Print("handlers::SaveURLJson -- write in ResponseWriter")

		hk := "content-type"
		hv := "application/json"
		w.Header().Set(hk, hv)
		log.Printf("handlers::SaveURLJson -- write in ResponseWriter, header key: %s, header value: %s", hk, hv)

		w.WriteHeader(http.StatusCreated)
		log.Printf("handlers::SaveURLJson -- write in ResponseWriter, status: %v", http.StatusCreated)

		w.Write(respBody)
		log.Printf("handlers::SaveURLJson -- write in ResponseWriter, body: %s", string(respBody))
	}
}
