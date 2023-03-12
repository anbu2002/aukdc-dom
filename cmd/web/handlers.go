package main

import (
//	"errors"
	"net/http"
//	"strconv"
//	"fmt"
	"time"

//	"aukdc.dom.com/internal/models"
	"aukdc.dom.com/internal/validator"

//	"github.com/julienschmidt/httprouter"
)
type facultySignupForm struct{
	FacultyID int `form:"facid"`
	Name string `form:"name"`
	PhoneNumber int64 `form:"phnno"`
	Email string `form:"email"`
	FacultyType string `form:"factype"`
	Department string `form:"dep"`
	Designation string `form:"designation"`
	Password []byte `form:"passd"`
	PanID string `form:"panid"`
//	PanPicture
	ExtensionNumber int64 `form:"extnumber"`
//	Esign 
	validator.Validator `form:"-"`
}

type facultyLoginForm struct{
	FacultyID string `form:"facid"`
	Password string `form:"pass"`
	validator.Validator `form:"-"`
}

type honorariumCreateForm struct {
	TransactionID int `form:"transid"`
	FacultyID string `form:"facid"`
	CourseCode string `form:"coursecode"`
//is this user input?
	InitialAmount int `form:"-"`
	FinalAmount int `form:"-"`
	CreatedTime time.Time `form:"-"`
}
type qPKCreateForm struct {
	honorariumCreateForm
	QuestionPaperCount int `form:"qpcount"`
	KeyCount int `form:"keycount"`
//is this user input?
	QuestionPaperRate float32 `form:"-"`
	KeyRate float32 `form:"-"`
}
type valuedPaper struct {
	honorariumCreateForm
	AnswerScriptCount int `form:"ascount"`
//is this user input?
	AnswerScriptRate float32 `form:"-"`
}
//Method Stub
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	
}

//Method Stub
func (app *application) honorariumView(w http.ResponseWriter, r *http.Request) {
}

//Method Stub
func (app *application) honorariumCreate(w http.ResponseWriter, r *http.Request) {
}

//Method Stub
func (app *application) honorariumCreatePost (w http.ResponseWriter, r *http.Request) {
}

//Method Stub
func (app *application) facultySignup(w http.ResponseWriter, r *http.Request) {
}

//Method Stub
func (app *application) facultySignupPost(w http.ResponseWriter, r *http.Request) {
}

//Method Stub
func (app *application) facultyLogin(w http.ResponseWriter, r *http.Request) {

}

//Method Stub
func (app *application) facultyLoginPost(w http.ResponseWriter, r *http.Request) {
}

//Method Stub
func (app *application) facultyLogoutPost(w http.ResponseWriter, r *http.Request) {
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
