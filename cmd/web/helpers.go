package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"io"
	"mime/multipart"
	"os"
	"strings"
	"strconv"

	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {

	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
	}
	buf := new(bytes.Buffer)
	var err error
	if strings.Contains(page,"print"){
		err=ts.ExecuteTemplate(buf,"printb",data)
	}else {
		err= ts.ExecuteTemplate(buf, "base", data)
	}
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}
func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		IsAuthorized:    app.isAuthorized(r),
		CSRFToken:       nosurf.Token(r),
	}
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		return err
	}
	return err
}
func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}
func (app *application) isAuthorized(r *http.Request) bool {
	isAuthorized, ok := r.Context().Value(isAuthorizedContextKey).(bool)
	if !ok {
		return false
	}
	return isAuthorized
}
func (app *application) isFaculty(r *http.Request) bool {
	isFaculty, ok := r.Context().Value(isFacultyContextKey).(bool)
	if !ok {
		return false
	}
	return isFaculty
}
func (app *application) hasBankDetails(r *http.Request) bool {
	hasBankDetails, ok := r.Context().Value(hasBankDetailsContextKey).(bool)
	if !ok {
		return false
	}
	return hasBankDetails
}
//to optimize
func (app *application) uploadImage(picture multipart.File, picID, imageName string) (error) {

	file, err := os.OpenFile("uploads/"+picID+"/"+imageName, os.O_EXCL|os.O_WRONLY|os.O_CREATE, 0400)

	defer file.Close()
	if err != nil {
		return err
	}
	io.Copy(file, picture)

	return nil
}
func (app *application) createImagePath(handler *multipart.FileHeader, picID string, id int)(string){

	facID:=strconv.Itoa(id)
	splitsName := strings.Split(handler.Filename, ".")
	var lenfilename = len(splitsName)

	image:= facID+"_"+picID+"."+splitsName[lenfilename-1]

	return image
}
