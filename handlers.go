package main

import (
	"html/template"
	"net/http"

	"github.com/dchest/captcha"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseGlob("templates/*.html"))
}

func indexPageHandler(w http.ResponseWriter, req *http.Request) {
	err := tmpl.ExecuteTemplate(w, "index", nil)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func loginPageHandler(w http.ResponseWriter, req *http.Request) {
	// session, _ := store.Get(req, "user")
	if req.Method == http.MethodGet {

		// if session.Values["authenticated"] != nil || session.Values["authenticated"] != false {
		// 	if session.Values["role"].(string) == "moderator" {

		// 		http.Redirect(w, req, "/complaints", 301)
		// 		return
		// 	}
		// 	if session.Values["role"].(string) == "member" {
		// 		http.Redirect(w, req, "/home", 301)
		// 		return
		// 	}
		// }
		captchaID := captcha.NewLen(6)
		err := tmpl.ExecuteTemplate(w, "login", struct {
			CSRF    template.HTML
			CAPTCHA string
		}{
			CSRF:    csrf.TemplateField(req),
			CAPTCHA: captchaID,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if req.Method == http.MethodPost {
		login(w, req)
		return
	}

}

func createPageHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		session, _ := store.Get(req, "user")
		data := struct {
			Email    string
			Fullname string
			CSRF     template.HTML
		}{Email: session.Values["email"].(string), Fullname: session.Values["fullname"].(string), CSRF: csrf.TemplateField(req)}
		err := tmpl.ExecuteTemplate(w, "create", data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if req.Method == http.MethodPost {
		create(w, req)
		return
	}

}

func updatePageHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	if req.Method == http.MethodGet {
		Complaint := dbComplaintByID(vars["id"])
		if Complaint == nil {
			http.Error(w, "Could not load complaint", http.StatusInternalServerError)
			return
		}
		data := struct {
			Complaint complaint
			CSRF      template.HTML
		}{Complaint: *Complaint, CSRF: csrf.TemplateField(req)}
		err := tmpl.ExecuteTemplate(w, "update", data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if req.Method == http.MethodPost {
		update(vars["id"], w, req)
		return
	}

}

func signupPageHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		captchaID := captcha.NewLen(6)
		err := tmpl.ExecuteTemplate(w, "signup", struct {
			CSRF    template.HTML
			CAPTCHA string
		}{
			CSRF:    csrf.TemplateField(req),
			CAPTCHA: captchaID,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if req.Method == http.MethodPost {
		signup(w, req)
		return
	}
}
func homePageHandler(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "user")
	Complaints, err := dbComplaintsByEmail(session.Values["email"].(string))
	if err != nil {
		http.Error(w, "Could not load complaints", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "home", struct {
		CSRF       template.HTML
		Complaints []complaint
	}{
		CSRF:       csrf.TemplateField(req),
		Complaints: Complaints,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func complaintsPageHandler(w http.ResponseWriter, req *http.Request) {
	Complaints, err := dbComplaintsAll()
	if err != nil {
		http.Error(w, "Could not load complaints", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "complaints", struct {
		CSRF       template.HTML
		Complaints []complaint
	}{
		CSRF:       csrf.TemplateField(req),
		Complaints: Complaints,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func membersPageHandler(w http.ResponseWriter, req *http.Request) {
	Members, err := dbMembersAll()
	if err != nil {
		http.Error(w, "Could not load members", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "members", struct {
		CSRF    template.HTML
		Members []member
	}{
		CSRF:    csrf.TemplateField(req),
		Members: Members,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
