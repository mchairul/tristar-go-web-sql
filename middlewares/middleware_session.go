package middlewares

import (
	"net/http"
	"websql/helpers"
)

func SessionMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := helpers.GetSessionStore(r)

		if err != nil {
			// jika gagal mendapatkan session
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if auth, ok := session.Values["Authenticated"].(bool); !ok || !auth {
			// jika tidak terautentikasi
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		//jika berhasil terauntentikasi
		next.ServeHTTP(w, r)
	}
}
