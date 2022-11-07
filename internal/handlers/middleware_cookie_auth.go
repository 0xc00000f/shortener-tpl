package handlers

import (
	"net/http"
	"time"

	"github.com/0xc00000f/shortener-tpl/internal/user"
)

const CookieAuthName = "Authorization"

func CookieAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		expires := time.Now().AddDate(0, 1, 0)
		reqCk, err := req.Cookie(CookieAuthName)
		if err != nil && err != http.ErrNoCookie {
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}

		if err == nil && user.Valid(reqCk.Value) {
			reqCk.Expires = expires
			http.SetCookie(w, reqCk)
			next.ServeHTTP(w, req)
			return
		}

		ck := http.Cookie{
			Name:    CookieAuthName,
			Path:    "/",
			Expires: expires,
		}

		u := user.New()
		eu, err := u.UserEncryptToString()
		if err != nil {
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}

		ck.Value = eu

		req.AddCookie(&ck)
		http.SetCookie(w, &ck)
		next.ServeHTTP(w, req)
	})
}

func GetUserFromRequest(r *http.Request) (u user.User, ok bool) {
	ck, err := r.Cookie(CookieAuthName)
	if err != nil {
		return user.User{}, false
	}

	err = u.UserDecryptFromString(ck.Value)
	if err != nil {
		return user.User{}, false
	}
	return u, true
}
