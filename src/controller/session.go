/*
 * Flashlight - By Asef Alper Tunga Dündar
 *
 * This Appication is something like Instagram for the University of Applied Sciences
 *
 * API version: 1.0.0
 */

package controller

import (
	"net/http"
	"utils"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func init() {
	key := make([]byte, 32)
	//rand.Read(key)
	key = []byte("I'm there for convinience!")
	store = sessions.NewCookieStore(key)
}

// Auth is used on all Handlers after a specitfic Client has logged in on Flashlight. It is a
// middleare and is to be chained the repsective handerls after passing the authentification
// process it will redirect the Client to our Mainpage and sets the Authetification Cookie.
// With this Cookie the Frontend Applications registers the Cookie and pushes the Client
// "/users/{userid}" Page.
// Returns HTTP-Status Code 401 Unauthorized when authorization failed
// Redirects to respective handler when authorization succeeded
func Auth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "session")

		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
		}

		// Check if user is authenticated
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Error(w, "You are not logged in", http.StatusUnauthorized)
		} else {

			h(w, r)
			utils.LoggerCookie(session)

		}
	}
}
