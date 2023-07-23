package main

import (
	"net/http"

	"aukdc.dom.com/ui"
	
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.FS(ui.Files))

	router.Handler(http.MethodGet, "/static/*filepath",fileServer)

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
	
	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))
	router.Handler(http.MethodGet, "/honorarium/view/:tyid/:hid", protected.ThenFunc(app.honorariumView))
	router.Handler(http.MethodGet, "/honorarium/view/:tyid/:hid/generate", protected.ThenFunc(app.generatePrint))
	

	faculty:=protected.Append(app.requireFaculty, app.checkFacultyDetails)
	router.Handler(http.MethodGet, "/faculty/details", faculty.ThenFunc(app.addDetails))
	router.Handler(http.MethodPost, "/faculty/details", faculty.ThenFunc(app.addDetailsPost))

	faculty=faculty.Append(app.requireFacultyDetails)
	router.Handler(http.MethodGet, "/honorarium/view/", faculty.ThenFunc(app.facultyView))
	router.Handler(http.MethodGet, "/honorarium/create/qpk", faculty.ThenFunc(app.qpkCreate))
	router.Handler(http.MethodPost, "/honorarium/create/qpk", faculty.ThenFunc(app.qpkCreatePost))
	router.Handler(http.MethodGet, "/honorarium/create/ansv", faculty.ThenFunc(app.ansvCreate))
	router.Handler(http.MethodPost, "/honorarium/create/ansv", faculty.ThenFunc(app.ansvCreatePost))

	authorized := protected.Append(app.requireAuthority)

	router.Handler(http.MethodGet, "/honorarium/", authorized.ThenFunc(app.facultyView))
	router.Handler(http.MethodGet, "/faculty/view", authorized.ThenFunc(app.facultyViewAll))
	router.Handler(http.MethodGet, "/faculty/view/:fid", authorized.ThenFunc(app.facultyView))
	router.Handler(http.MethodGet, "/faculty/view/:fid/honorarium/:hid/print", authorized.ThenFunc(app.generatePrint))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
}
