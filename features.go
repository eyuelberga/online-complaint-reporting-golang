package main

import (
	"html"
	"net/http"
)

func create(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "user")
	filename, stop := upload(w, r)
	if stop {
		return
	}

	created := dbCreateComplaint(complaintPayload{email: session.Values["email"].(string), fullname: session.Values["fullname"].(string), comment: html.EscapeString(r.FormValue("comment")), file: filename})
	if created {
		http.Redirect(w, r, "/home", 301)
		return
	}
	http.Error(w, "Could not Create Complaint", http.StatusInternalServerError)
	return

}

func update(id string, w http.ResponseWriter, r *http.Request) {
	filename, stop := upload(w, r)
	if stop {
		return
	}
	if filename == "" {
		created := dbUpdateComplaintComment(id, struct {
			comment string
		}{comment: html.EscapeString(r.FormValue("comment"))})
		if created {
			http.Redirect(w, r, "/home", 301)
			return
		}
		return
	}

	created := dbUpdateComplaint(id, struct {
		comment string
		file    string
	}{comment: html.EscapeString(r.FormValue("comment")), file: filename})
	if created {
		http.Redirect(w, r, "/home", 301)
		return
	}

	http.Error(w, "Could not Update Complaint", http.StatusInternalServerError)
	return

}

func enableAccount(w http.ResponseWriter, r *http.Request) {
	updated := dbUpdateMemberAccount(r.FormValue("email"), true)
	if updated {
		http.Redirect(w, r, "/members", 301)
		return
	}

	http.Error(w, "Could not Update Member Account", http.StatusInternalServerError)
	return

}

func disableAccount(w http.ResponseWriter, r *http.Request) {
	updated := dbUpdateMemberAccount(r.FormValue("email"), false)
	if updated {
		http.Redirect(w, r, "/members", 301)
		return
	}

	http.Error(w, "Could not Update Member Account", http.StatusInternalServerError)
	return

}
