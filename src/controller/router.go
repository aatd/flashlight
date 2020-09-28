/*
 * Flashlight - By Asef Alper Tunga Dündar
 *
 * This Appication is something like Instagram for the University of Applied Sciences
 *
 * API version: 1.0.0
 */

package controller

import (
	"html/template"
	"net/http"
	"utils"

	"github.com/gorilla/mux"
)

var tmpl *template.Template

func init() {
	//Parse all views when starting the server
	tmpl = template.Must(template.ParseGlob("views/*.html"))
}

// Routes the Collection of all Routes for Flashlight
type Routes []Route

// Route defines all functios and metadata of a Route itself
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// CreateRouter Creates a new Router for the Flashlight Application. Reads all
// Routes from the controller.Routes and handles the routing for the static content
// (css/js/icons) of Flashlight
func CreateRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)

	//Static Files for everyone
	router.PathPrefix("/icon/").Handler(http.StripPrefix("/icon", http.FileServer(http.Dir("./src/static/icon"))))
	router.PathPrefix("/css/").Handler(http.StripPrefix("/css", http.FileServer(http.Dir("./src/static/css"))))
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js", http.FileServer(http.Dir("./src/static/js"))))

	//Static Files when logged in
	router.PathPrefix("/users/icon/").Handler(http.StripPrefix("/users/icon", http.FileServer(http.Dir("./src/static/icon"))))
	router.PathPrefix("/users/css/").Handler(http.StripPrefix("/users/css", http.FileServer(http.Dir("./src/static/css"))))
	router.PathPrefix("/users/js/").Handler(http.StripPrefix("/users/js", http.FileServer(http.Dir("./src/static/js"))))

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = utils.Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	//Handle Queries
	router.Get("Login           ").Queries("action", "login")
	router.Get("Logout          ").Queries("action", "logout")
	router.Get("Register        ").Queries("action", "register")
	router.Get("Userdata        ").Queries("action", "userdata")
	router.Get("GetImages       ").Queries("lastRecordTime", "{[0-9]+}")

	return router
}

// Collection of all Routes for the Application. Middleware chaines are also defined wihting this collection
var routes = Routes{

	//////////////////////////////////////////////////////////////////////
	//							Public Routes							//
	//////////////////////////////////////////////////////////////////////

	Route{"GetMainpage     ",
		"GET",
		"/",
		GetMainPage,
	},

	Route{"GetImageMetaData",
		"GET",
		"/images/{imageID:[a-f0-9]+}",
		GetImageMetaData,
	},

	Route{"GetImage        ",
		"GET",
		"/images/{imageID:[a-f0-9]+}/raw",
		GetImage,
	},

	Route{"GetImages       ",
		"GET",
		"/images",
		GetImages,
	},

	Route{"Register        ",
		"POST",
		"/users",
		Register,
	},

	Route{"Login           ",
		"POST",
		"/users",
		Login,
	},

	//////////////////////////////////////////////////////////////////////
	//							Private Routes							//
	//////////////////////////////////////////////////////////////////////

	Route{"ProfilePage     ",
		"GET",
		"/users/{userID:[a-f0-9]+}",
		Auth(ProfilePage),
	},

	Route{"Userdata        ",
		"GET",
		"/users",
		Auth(Userdata),
	},

	Route{"Logout          ",
		"POST",
		"/users/{userID:[a-f0-9]+}",
		Auth(Logout),
	},

	Route{"GetUserImages   ",
		"GET",
		"/users/{userID:[a-f0-9]+}/images",
		Auth(GetUserImages),
	},

	Route{"GetLike         ",
		"GET",
		"/images/{imageID:[a-f0-9]+}/like",
		Auth(GetLike),
	},

	Route{"UploadImage     ",
		"POST",
		"/users/{userID:[a-f0-9]+}/images",
		Auth(UploadImage),
	},

	Route{"CommentImage    ",
		"POST",
		"/images/{imageID:[a-f0-9]+}/comment",
		Auth(CommentImage),
	},

	Route{"LikeImage       ",
		"POST",
		"/images/{imageID:[a-f0-9]+}/like",
		Auth(LikeImage),
	},

	Route{"DeleteImage     ",
		"DELETE",
		"/images/{imageID:[a-f0-9]+}",
		Auth(DeleteImage),
	},
}
