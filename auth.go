package main

import "net/http"

func signup(w http.ResponseWriter, r *http.Request) {

	userAdded := dbAddUser(registrationPayload{email: r.FormValue("email"), fullname: r.FormValue("fullname"), password: r.FormValue("password"), role: "member", enabled: true})
	if userAdded {
		http.Redirect(w, r, "/login?success=registration%20completed", 301)
		return
	}
	http.Error(w, "Could not Register", http.StatusInternalServerError)

}

func login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "user")
	user := dbGetUser(r.FormValue("email"))
	if user != nil {
		// Set user as authenticated
		authenticated := checkPasswordHash(r.FormValue("password"), user.password)
		if authenticated {
			if !user.enabled {
				http.Error(w, "Your account has been disabled", http.StatusForbidden)
				return
			}
			session.Values["authenticated"] = true
			session.Values["role"] = user.role
			session.Values["email"] = user.email
			session.Values["fullname"] = user.fullname
			session.Save(r, w)
			if user.role == "member" {
				http.Redirect(w, r, "/home", 301)
				return
			}

			if user.role == "moderator" {
				http.Redirect(w, r, "/complaints", 301)
				return
			}
		}

		http.Redirect(w, r, "/login?error=wrong%20password", 301)
		return
	}
	http.Error(w, "Could not Login", http.StatusInternalServerError)

}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "user")

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/", 301)
}
