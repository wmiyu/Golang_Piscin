/*
 * Team01 Server
 *
 * API version: 1.0.0
 */
package teamnode01

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"Connect",
		strings.ToUpper("Get"),
		"/connect",
		GetConnect,
	},
	Route{
		"SetBook",
		strings.ToUpper("Post"),
		"/setbook",
		SetBook,
	},
	Route{
		"GetBook",
		strings.ToUpper("Post"),
		"/getbook",
		GetBook,
	},
	Route{
		"DeleteBook",
		strings.ToUpper("Post"),
		"/deletebook",
		DeleteBook,
	},
	Route{
		"RegisterNode",
		strings.ToUpper("Post"),
		"/registernode",
		RegisterNode,
	},
}
