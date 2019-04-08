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

	dash := router.PathPrefix("/auth").Subrouter()
	dash.HandleFunc("/signin", loginHandler).Methods("POST")

	return router
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/health",
		index,
	},
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
		"UserAdd",
		"POST",
		"/users",
		userAdd,
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
		"RoleFindAll",
		"GET",
		"/roles",
		roleFindAll,
	},
	Route{
		"RoleFindOne",
		"GET",
		"/roles/{ID}",
		roleFindOne,
	},
	Route{
		"RoleAdd",
		"POST",
		"/roles",
		roleAdd,
	},
	Route{
		"RoleUpdate",
		"PUT",
		"/roles/{ID}",
		roleUpdate,
	},
	Route{
		"RoleDelete",
		"DELETE",
		"/roles/{ID}",
		roleDelete,
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
}
