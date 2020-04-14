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

// Auth Done
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
