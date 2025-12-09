package middlewares

import (
	"net/http"
	"websql/constants"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(constants.SessionScret))

func SessionMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, constants.SessionName)

		if err != nil {
			// jika gagal mendapatkan session
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
			return
		}

		userid, ok := session.Values["Userid"]
		if !ok || userid == nil {
			// jika tidak terautentikasi
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
			return
		}

		//jika berhasil terauntentikasi
		next.ServeHTTP(w, r)
	}
}
