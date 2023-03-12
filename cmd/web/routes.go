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
	router.Handler(http.MethodGet, "/honorarium/view/:id", dynamic.ThenFunc(app.honorariumView))
	router.Handler(http.MethodGet, "/faculty/signup", dynamic.ThenFunc(app.facultySignup))
	router.Handler(http.MethodPost, "/faculty/signup", dynamic.ThenFunc(app.facultySignupPost))
	router.Handler(http.MethodGet, "/faculty/login", dynamic.ThenFunc(app.facultyLogin))
	router.Handler(http.MethodPost, "/faculty/login", dynamic.ThenFunc(app.facultyLoginPost))

	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/honorarium/create", protected.ThenFunc(app.honorariumCreate))
	router.Handler(http.MethodPost, "/honorarium/create", protected.ThenFunc(app.honorariumCreatePost))
	router.Handler(http.MethodPost, "/faculty/logout", protected.ThenFunc(app.facultyLogoutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
}
