package main

import (
	"log"
	"net/http"

	rl "github.com/ahmedash95/ratelimit"
	"github.com/dchest/captcha"
	"github.com/gorilla/mux"
)

func logger(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		f(w, r)
	}
}

func validateCAPTCHA(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			id := r.FormValue("captcha_id")
			value := r.FormValue("captcha_solution")
			log.Println(id, value)
			if captcha.VerifyString(id, value) {
				f(w, r)
			} else {
				http.Redirect(w, r, r.URL.Path+"?error=CAPTCHA%20did%20not%20match", 301)
			}
		} else {
			f(w, r)
		}
	}
}

func hasRole(f http.HandlerFunc, role string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "user")
		// Check if user is authenticated
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			log.Println("[hasRole Middleware] not authenticated")
			http.Redirect(w, r, "/login?error=not%20authenticated", 301)
			return
		}
		if session.Values["role"] != role {
			log.Printf("[hasRole Middleware] expected role %s, but found %s", role, session.Values["role"])
			http.Redirect(w, r, "/login?error=not%20authorized", 301)
			return
		}
		f(w, r)
	}
}

func checkDownloadable(f http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		session, _ := store.Get(r, "user")
		// Check if user is authenticated
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			log.Println("[checkDownloadable Middleware] not authenticated")
			http.Redirect(w, r, "/login?error=not%20authenticated", 301)
			return
		}
		if session.Values["role"] != "moderator" {
			// if user is not a moderator, check if is owner
			isOwner := dbUserOwnsFile(session.Values["email"].(string), vars["file"])
			if !isOwner {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

		}
		f(w, r)
	}
}

func rateLimitter(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := "127.0.0.1"
		if !isValidRequest(ratelimit, ip) {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		ratelimit.Hit(ip)
		h(w, r)
	}
}

func isValidRequest(l rl.Limit, key string) bool {
	_, ok := l.Rates[key]
	if !ok {
		return true
	}
	if l.Rates[key].Hits == l.MaxRequests {
		return false
	}
	return true
}
