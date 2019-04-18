package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/metiago/zbx1/common/middleware"
	"github.com/urfave/negroni"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {

	apiRouter := mux.NewRouter().StrictSlash(false)

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = middleware.Token(middleware.Logger(handler, route.Name))
		apiRouter.PathPrefix("/api/v1").Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(handler)
	}

	router := mux.NewRouter()
	router.PathPrefix("/api").Handler(negroni.New(negroni.HandlerFunc(authHandler), negroni.Wrap(apiRouter)))

	// Public routes
	auth := router.PathPrefix("/auth").Subrouter()
	authHandler := http.HandlerFunc(loginHandler)
	auth.Handle("/login", middleware.Logger(authHandler, "/auth")).Methods("POST")

	public := router.PathPrefix("/").Subrouter()
	indexHandler := http.HandlerFunc(index)
	signUpHandler := http.HandlerFunc(userAdd)
	public.Handle("/", middleware.Logger(indexHandler, "/")).Methods("GET")
	public.Handle("/signup", middleware.Logger(signUpHandler, "/signup")).Methods("POST")

	return router
}

var routes = Routes{
	Route{
		"UserFindAll",
		"GET",
		"/users",
		userFindAll,
	},
	Route{
		"UserFindOne",
		"GET",
		"/users/{ID}",
		userFindOne,
	},
	Route{
		"UserUpdate",
		"PUT",
		"/users/{ID}",
		userUpdate,
	},
	Route{
		"UserDelete",
		"DELETE",
		"/users/{ID}",
		userDelete,
	},
	Route{
		"FileUpload",
		"POST",
		"/files/upload",
		fileUpload,
	},
	Route{
		"FileDownload",
		"GET",
		"/files/download/{ID}",
		fileDownload,
	},
	Route{
		"FileFindAllByUsername",
		"GET",
		"/files/{username}",
		fileFindAllByUsername,
	},
	Route{
		"FileDelete",
		"DELETE",
		"/files/{ID}",
		fileDelete,
	},
}
