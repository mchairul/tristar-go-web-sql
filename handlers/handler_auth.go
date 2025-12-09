package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"
	"websql/typecustom"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func HandleFormLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		username := r.FormValue("username")
		password := r.FormValue("password")
		fmt.Println(username, password)

		var result = typecustom.User{}

		err := db.QueryRow("SELECT * FROM users WHERE username = ?", username).
			Scan(&result.Id, &result.Username, &result.Password, &result.Name)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("no user with username " + username)
				http.Redirect(w, r, "/", http.StatusMovedPermanently)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		fmt.Println(result)

		err = bcrypt.CompareHashAndPassword(
			[]byte(result.Password), []byte(password))

		if err != nil {
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
			return
		}

		var store = sessions.NewCookieStore([]byte("secret-key"))
		session, _ := store.Get(r, "session-web")

		session.Values["userid"] = result.Id
		session.Values["username"] = result.Username
		session.Values["name"] = result.Name

		session.Save(r, w)

		w.Write(([]byte("Login Sukses")))
	}
}
