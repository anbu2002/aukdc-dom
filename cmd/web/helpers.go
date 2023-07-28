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
	"strings"

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
	if strings.Contains(page, "print") {
		err = ts.ExecuteTemplate(buf, "printb", data)
	} else {
		err = ts.ExecuteTemplate(buf, "base", data)
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

func (app *application) convertImage(file multipart.File)([]byte, error) {
	imageData := make([]byte, 0)
	buf := make([]byte, 4096)
	for {
		n, err := file.Read(buf)
		if n == 0 || err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		imageData = append(imageData, buf[:n]...)
	}
	return imageData, nil
}

