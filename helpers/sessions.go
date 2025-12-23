package helpers

import (
	"net/http"
	"websql/constants"

	"github.com/gorilla/sessions"
)

var Session = sessions.NewFilesystemStore(constants.SessionDirectory, []byte("secret-key"))

func GetSessionStore(r *http.Request) (*sessions.Session, error) {
	// Session.Options = &sessions.Options{
	// 	Path:     "/",
	// 	MaxAge:   0,
	// 	HttpOnly: true,
	// }
	return Session.Get(r, constants.SessionName)
}
