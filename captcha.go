package main

import (
	"log"
	"net/http"

	"github.com/dchest/captcha"
	"github.com/gorilla/mux"
)

func captchaImage(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	if err := captcha.WriteImage(w, id, 120, 80); err != nil {
		log.Println("[captchaImage]", err)
	}
}
