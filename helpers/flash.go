package helpers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

func SetFlash(session *sessions.Session, r *http.Request, w http.ResponseWriter, name, value string) {
	session.AddFlash(value, name)
	session.Save(r, w)
}

func GetFlash(session *sessions.Session, r *http.Request, w http.ResponseWriter, name string) []string {
	flashMessage := session.Flashes(name)

	if len(flashMessage) > 0 {
		err := session.Save(r, w)
		if err != nil {
			fmt.Println(err.Error())
		}

		var flashes []string

		for _, flashList := range flashMessage {
			flashes = append(flashes, flashList.(string))
		}

		return flashes
	}
	return nil
}
