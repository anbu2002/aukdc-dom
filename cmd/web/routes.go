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
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
	
	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/", protected.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/honorarium/view/", protected.ThenFunc(app.facultyView))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))
	router.Handler(http.MethodGet, "/honorarium/view/:tyid/:hid", protected.ThenFunc(app.honorariumView))
	router.Handler(http.MethodGet, "/honorarium/view/:tyid/:hid/generate", protected.ThenFunc(app.generatePrint))
	

	faculty:=protected.Append(app.requireFaculty, app.checkBankDetails)
	router.Handler(http.MethodGet, "/faculty/bankdetails", faculty.ThenFunc(app.addBankDetails))
	router.Handler(http.MethodPost, "/faculty/bankdetails", faculty.ThenFunc(app.addBankDetailsPost))

	faculty=faculty.Append(app.requireBankDetails)
	router.Handler(http.MethodGet, "/honorarium/qpk/create", faculty.ThenFunc(app.qpkCreate))
	router.Handler(http.MethodPost, "/honorarium/qpk/create", faculty.ThenFunc(app.qpkCreatePost))
	router.Handler(http.MethodGet, "/honorarium/ansv/create", faculty.ThenFunc(app.ansvCreate))
	router.Handler(http.MethodPost, "/honorarium/ansv/create", faculty.ThenFunc(app.ansvCreatePost))

	authorized := protected.Append(app.requireAuthority)

	router.Handler(http.MethodGet, "/faculty/view", authorized.ThenFunc(app.facultyViewAll))
	router.Handler(http.MethodGet, "/faculty/view/:fid", authorized.ThenFunc(app.facultyView))
	router.Handler(http.MethodGet, "/faculty/view/:fid/honorarium/:hid/print", authorized.ThenFunc(app.generatePrint))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
}
