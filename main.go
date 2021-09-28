package main

import (
	"log"
	"net/http"
	"os"
	"time"

	rl "github.com/ahmedash95/ratelimit"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var ratelimit rl.Limit

func main() {
	godotenv.Load()
	db = newDatabase(os.Getenv("DATABASE_URL"))
	r := mux.NewRouter()
	r.Use(mux.CORSMethodMiddleware(r))
	r.HandleFunc("/", logger(indexPageHandler)).Methods("GET")
	r.HandleFunc("/home", hasRole(logger(homePageHandler), "member")).Methods("GET")
	r.HandleFunc("/complaints", hasRole(logger(complaintsPageHandler), "moderator")).Methods("GET")
	r.HandleFunc("/members", hasRole(logger(membersPageHandler), "moderator")).Methods("GET")
	r.HandleFunc("/disable", hasRole(logger(disableAccount), "moderator")).Methods("POST")
	r.HandleFunc("/enable", hasRole(logger(enableAccount), "moderator")).Methods("POST")
	r.HandleFunc("/uploads/{file}", checkDownloadable(logger(downloadHandler))).Methods("GET")
	r.HandleFunc("/captcha/{id}", logger(captchaImage)).Methods("GET")
	r.HandleFunc("/login", validateCAPTCHA(logger(loginPageHandler))).Methods("GET", "POST")
	r.HandleFunc("/logout", logger(logout)).Methods("POST")
	r.HandleFunc("/new", hasRole(logger(createPageHandler), "member")).Methods("GET", "POST")
	r.HandleFunc("/update/{id}", hasRole(logger(updatePageHandler), "member")).Methods("GET", "POST")
	r.HandleFunc("/signup", validateCAPTCHA(logger(signupPageHandler))).Methods("GET", "POST")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	srv := &http.Server{
		Handler:      csrf.Protect([]byte(os.Getenv("CSRF_KEY")))(r),
		Addr:         "0.0.0.0:9090",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  time.Second * 60,
	}

	log.Fatal(srv.ListenAndServe())
}
