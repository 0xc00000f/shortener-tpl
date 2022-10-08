package helpers

import "net/http"

func BadRequest(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "400 page not found", http.StatusBadRequest)
}
