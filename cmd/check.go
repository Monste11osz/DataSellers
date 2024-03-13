package main

import "net/http"

func RequireCookie(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(COOKIE)
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/us/authentication", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}
