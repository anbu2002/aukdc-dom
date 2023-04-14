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

	router.Handler(http.MethodGet, "/faculty/signup", dynamic.ThenFunc(app.facultySignup))
	router.Handler(http.MethodPost, "/faculty/signup", dynamic.ThenFunc(app.facultySignupPost))
	router.Handler(http.MethodGet, "/faculty/login", dynamic.ThenFunc(app.facultyLogin))
	router.Handler(http.MethodPost, "/faculty/login", dynamic.ThenFunc(app.facultyLoginPost))

	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/", protected.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/honorarium/view/:id", protected.ThenFunc(app.honorariumView))
	router.Handler(http.MethodGet, "/faculty/bankdetails", protected.ThenFunc(app.addBankDetails))
	router.Handler(http.MethodPost, "/faculty/bankdetails", protected.ThenFunc(app.addBankDetailsPost))
	router.Handler(http.MethodGet, "/honorarium/qpk/create", protected.ThenFunc(app.qpkCreate))
	router.Handler(http.MethodPost, "/honorarium/qpk/create", protected.ThenFunc(app.qpkCreatePost))
	router.Handler(http.MethodGet, "/honorarium/ansv/create", protected.ThenFunc(app.ansvCreate))
	router.Handler(http.MethodPost, "/honorarium/ansv/create", protected.ThenFunc(app.ansvCreatePost))
	router.Handler(http.MethodPost, "/faculty/logout", protected.ThenFunc(app.facultyLogoutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
}
