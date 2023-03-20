package main

import (
	"errors"
	"net/http"
//	"strconv"
	"fmt"
	"time"

	"aukdc.dom.com/internal/models"
	"aukdc.dom.com/internal/validator"

//	"github.com/julienschmidt/httprouter"
)

//Please change the data types of picture variables.
type facultySignupForm struct{
	FacultyID string `form:"facultyid"`
	Name string `form:"name"`
	Phone int64 `form:"phone"`
	Email string `form:"email"`
	FacultyType string `form:"facultytype"`
	Department string `form:"dept"`
	Designation string `form:"designation"`
	Password string `form:"password"`
	PanID string `form:"panid"`
	PanPicture any `form:"panpic"`
	Extension int64 `form:"extnumber"`
	Esign any `form:"esign"`
	validator.Validator `form:"-"`
}

type facultyLoginForm struct{
	FacultyID string `form:"facultyid"`
	Password string `form:"password"`
	validator.Validator `form:"-"`
}

type honorariumCreateForm struct {
	TransactionID int `form:"transactionid"`
	FacultyID string `form:"facultyid"`
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
	honoraria,err:=app.honoraria.ViewAll()
	if err!=nil{
		app.serverError(w,err)
		return
	}
	data:=app.newTemplateData(r)
	data.Honoraria=honoraria

	app.render(w, http.StatusOK, "home.tmpl",data)
}

func (app *application) facultySignup(w http.ResponseWriter, r *http.Request) {
	if(app.isAuthenticated(r)){
		http.Redirect(w,r,"/",http.StatusSeeOther)
	}
	data:=app.newTemplateData(r)
	data.Form=facultySignupForm{}
	app.render(w, http.StatusOK,"signup.tmpl",data)

}

func (app *application) facultyLogin(w http.ResponseWriter, r *http.Request) {
	if(app.isAuthenticated(r)){
		http.Redirect(w,r,"/",http.StatusSeeOther)
	}
	data:=app.newTemplateData(r)
	data.Form=facultyLoginForm{}
	app.render(w, http.StatusOK,"login.tmpl",data)
}

func (app *application) facultySignupPost(w http.ResponseWriter, r *http.Request) {
	var form facultySignupForm

	err:=app.decodePostForm(r, &form)
	if err!=nil{
		app.clientError(w, http.StatusBadRequest)
		return
	}
        form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
        form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
        form.CheckField(validator.NotBlank(form.FacultyID), "facultyid", "This field cannot be blank")
        form.CheckField(validator.NotBlank(form.Department), "department", "This field cannot be blank")
        form.CheckField(validator.IntegerRange(form.Phone, 6000000000, 9999999999), "phone", "This field must be a valid phone number")
        form.CheckField(validator.IntegerRange(form.Extension, 6000000000, 9999999999), "phone", "This field must be a valid extension number")
        form.CheckField(validator.NotBlank(form.Designation), "designation", "This field cannot be blank")
        form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
        form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
        form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid(){
		data:=app.newTemplateData(r)
		data.Form=form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}
	
	err=app.faculty.Insert(form.FacultyID, form.Name, form.Phone, form.Email, form.FacultyType, form.Department, form.Designation, form.Password, form.PanID, "Insert picture", form.Extension, "Insert Picture")

	if err!=nil{
		//add error for duplicate id
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "You have sucessfully signed up. Please log in.")
	http.Redirect(w, r, "/faculty/login",http.StatusSeeOther)
}
func (app *application) facultyLoginPost(w http.ResponseWriter, r *http.Request) {
	var form facultyLoginForm

	err:=app.decodePostForm(r, &form)
	if err!=nil{
		app.clientError(w, http.StatusBadRequest)
		return
	}
        form.CheckField(validator.NotBlank(form.FacultyID), "facultyid", "This field cannot be blank")
        form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
        if !form.Valid() {
                data := app.newTemplateData(r)
                data.Form = form
                app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
                return
        }
	id, err:=app.faculty.Authenticate(form.FacultyID, form.Password)
	if err!=nil{
		if errors.Is(err, models.ErrInvalidCredentials){
			form.AddNonFieldError("Faculty ID or password is incorrect")
			data:=app.newTemplateData(r)
			data.Form=form
			app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		} else{
			app.serverError(w,err)
		}
		return
	}
	err=app.sessionManager.RenewToken(r.Context())
	if err!=nil{
		app.serverError(w, err)
		return
	}
	fmt.Println(id)
	app.sessionManager.Put(r.Context(),"authenticatedFacultyID",id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) facultyLogoutPost(w http.ResponseWriter, r *http.Request) {
	err:=app.sessionManager.RenewToken(r.Context())
        if err!=nil{
                app.serverError(w, err)
                return
        }

        app.sessionManager.Remove(r.Context(), "authenticatedFacultyID")
        app.sessionManager.Put(r.Context(),"flash","You've been logged out successfully!")
        http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) honorariumView(w http.ResponseWriter, r *http.Request) {
}

//Method Stub
func (app *application) honorariumCreate(w http.ResponseWriter, r *http.Request) {
}

//Method Stub
func (app *application) honorariumCreatePost (w http.ResponseWriter, r *http.Request) {
}


//Method Stub

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
