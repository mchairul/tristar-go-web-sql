package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"
	"websql/helpers"
	"websql/typecustom"

	"golang.org/x/crypto/bcrypt"
)

func HandleFormLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		helpers.SetHeaders(w, r)
		var data = typecustom.WebData{
			"Title": "Login",
		}

		var tmpl = template.Must(template.ParseFiles(
			"views/_header.html",
			"views/_footer.html",
			"views/login.html",
		))

		err := tmpl.ExecuteTemplate(w, "login", data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func HandlePostLogin(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		helpers.SetHeaders(w, r)
		username := r.FormValue("username")
		password := r.FormValue("password")

		var result = typecustom.User{}

		err := db.QueryRow("SELECT * FROM users WHERE username = ?", username).
			Scan(&result.Id, &result.Username, &result.Password, &result.Name)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("no user with username " + username)
				http.Redirect(w, r, "/", http.StatusSeeOther)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		err = bcrypt.CompareHashAndPassword(
			[]byte(result.Password), []byte(password))

		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		session, _ := helpers.GetSessionStore(r)

		session.Values["Userid"] = result.Id
		session.Values["Username"] = result.Username
		session.Values["Name"] = result.Name
		session.Values["Authenticated"] = true

		err = session.Save(r, w)

		if err != nil {
			http.Error(w, "gagal Save Session "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/listkaryawan", http.StatusSeeOther)
	}
}

func HandleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("*** Logging Out! ***")
		helpers.SetHeaders(w, r)

		session, err := helpers.GetSessionStore(r)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println(session.Values["Username"])

		delete(session.Values, "Authenticated")
		delete(session.Values, "Userid")
		delete(session.Values, "Username")
		delete(session.Values, "Name")

		//session.Values["Authenticated"] = false

		err = session.Save(r, w)

		if err != nil {
			http.Error(w, "gagal logout", http.StatusInternalServerError)
			return
		} else {
			fmt.Println(session)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}
}
