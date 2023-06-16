package main

import(
	"context"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

func secureHeaders(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Security-Policy",
		"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w,r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		defer func(){
			if err:= recover(); err!=nil{
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w,r)
	})
}

func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:"/",
		Secure:true,
	})
	return csrfHandler
}	
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id:=app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
        	if id==0 {
			next.ServeHTTP(w, r)
               		return
                 }
                 exists,err:=app.user.Exists(id)
                 if err!=nil {
                        app.serverError(w, err)
                        return
                 }
                 if exists {
                        ctx:=context.WithValue(r.Context(), isAuthenticatedContextKey, true)
                        r=r.WithContext(ctx)
			id=app.sessionManager.GetInt(r.Context(), "authorizedUserID")
                        if id==0 {
                                ctx=context.WithValue(r.Context(), isFacultyContextKey, true)
                        	r=r.WithContext(ctx)
                                next.ServeHTTP(w, r)
                                return
                        }
                        authorized,err:=app.user.Authorized(id)
                        if err!=nil {
                                app.serverError(w, err)
                                return
                        }
                        if authorized{
                                ctx=context.WithValue(r.Context(), isAuthorizedContextKey, true)
                        	r=r.WithContext(ctx)
                        }
                 }
                 next.ServeHTTP(w, r)
        })
}
func (app *application) checkBankDetails(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id:=app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
		exists,err:=app.user.HasBankDetails(id)
		if err!=nil {
			app.serverError(w, err)
			return
		}
		if exists{
			ctx:=context.WithValue(r.Context(), hasBankDetailsContextKey, true)
			r=r.WithContext(ctx)
		}
                next.ServeHTTP(w, r)
	})
}
func (app *application) requireAuthentication(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		if !app.isAuthenticated(r){
                	app.sessionManager.Put(r.Context(),"flash","Please login to access this page!")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		w.Header().Add("Cache-Control","no-store")
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireFaculty(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		if !app.isFaculty(r){
                	app.sessionManager.Put(r.Context(),"flash","You need to be faculty to access this page!")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		w.Header().Add("Cache-Control","no-store")
		next.ServeHTTP(w, r)
	})
}
func (app *application) requireAuthority(next http.Handler) http.Handler{
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
                if !app.isAuthorized(r){
                        app.sessionManager.Put(r.Context(),"flash","You are not authorized to view this page!")
                        http.Redirect(w, r, "/", http.StatusSeeOther)
                        return
                }

                w.Header().Add("Cache-Control","no-store")
                next.ServeHTTP(w, r)
        })
}
func (app *application) requireBankDetails(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		if !app.hasBankDetails(r){
			fmt.Println("here")
			app.sessionManager.Put(r.Context(), "flash", "Please add your bank details")
			http.Redirect(w, r, "/faculty/bankdetails", http.StatusSeeOther)
			return
		} 
                w.Header().Add("Cache-Control","no-store")
                next.ServeHTTP(w, r)
	})
}
