package main

import (
	"errors"
	"net/http"
	"strconv"
	"fmt"
	"time"

	"aukdc.dom.com/internal/models"
	"aukdc.dom.com/internal/validator"

//	"github.com/julienschmidt/httprouter"
)

//Please change the data types of picture variables.
type facultySignupForm struct{
	FacultyID int `form:"facultyid"`
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

type bankDetailsForm struct{ 
	BankName string `form:"bankname"`
	AccountNumber int `form:"accountno"`
	IFSC string `form:"IFSC"`
	Passbook any `"form:passpic"`
	validator.Validator `form:"-"`
}
type facultyLoginForm struct{
	FacultyID int `form:"facultyid"`
	Password string `form:"password"`
	validator.Validator `form:"-"`
}

type honorariumCreateForm struct {
	TransactionID int `form:"-"`
	FacultyID int `form:"-"`
	CourseCode string `form:"coursecode"`
	InitialAmount int `form:"-"`
	FinalAmount int `form:"-"`
	CreatedTime time.Time `form:"-"`
	validator.Validator `form:"-"`
}
type qpkCreateForm struct {
	honorariumCreateForm
	QuestionPaperCount int `form:"qc"`
	KeyCount int `form:"kc"`
	QuestionPaperRate float32 `form:"-"`
	KeyRate float32 `form:"-"`
}
type ansvCreateForm struct {
	honorariumCreateForm
	AnswerScriptCount int `form:"asc"`
	AnswerScriptRate float32 `form:"-"`
}
//Method Stub
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	honoraria,err:=app.honorarium.ViewAll()
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
        form.CheckField(validator.NotBlank(form.PanID), "panid", "This field cannot be blank")
        form.CheckField(validator.NotBlank(strconv.Itoa(form.FacultyID)), "facultyid", "This field cannot be blank")
        form.CheckField(validator.NotBlank(form.Department), "dept", "This field cannot be blank")
        form.CheckField(validator.IntegerRange(form.Phone, 6000000000, 9999999999), "phone", "This field must be a valid phone number")
        form.CheckField(validator.IntegerRange(form.Extension, 6000000000, 9999999999), "extnumber", "This field must be a valid extension number")
        form.CheckField(validator.NotBlank(form.Designation), "designation", "This field cannot be blank")
        form.CheckField(validator.NotBlank(form.FacultyType), "facultytype", "This field cannot be blank")
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
func (app *application) addBankDetails(w http.ResponseWriter, r *http.Request){
	data:=app.newTemplateData(r)
	flag:=false
	id:=app.sessionManager.Get(r.Context(),"authenticatedFacultyID").(int)
	bd, err:=app.faculty.GetBankDetails(id)
	if err!=nil{
		if errors.Is(err, models.ErrNoRecord){
			flag=true
		} else{
			app.serverError(w,err)
			return
		}
	}
	if flag{
		data.Form = bankDetailsForm{}
	}else {
		data.Form = bankDetailsForm{
			BankName:bd.BankName,
			AccountNumber:bd.AccountNumber,
			IFSC:bd.IFSC,
		}
	}
	app.render(w, http.StatusOK, "bankdetails.tmpl", data)
}

func (app *application) addBankDetailsPost(w http.ResponseWriter, r *http.Request){
	id:=app.sessionManager.Get(r.Context(),"authenticatedFacultyID").(int)
	var form bankDetailsForm
	err:=app.decodePostForm(r, &form)
	if err!=nil{
		app.clientError(w, http.StatusBadRequest)
		return
	}
        form.CheckField(validator.NotBlank(form.BankName), "bank", "This field must not be blank")
        form.CheckField(validator.NotBlank(strconv.Itoa(form.AccountNumber)), "accountno", "This field must be a valid account number")
        form.CheckField(validator.MinChars(form.IFSC, 11), "IFSC", "IFSC code must be exactly 11 characters")
        form.CheckField(validator.MaxChars(form.IFSC, 11), "IFSC", "IFSC code must be exactly 11 characters")
	if !form.Valid(){
		data:=app.newTemplateData(r)
		data.Form=form
		app.render(w, http.StatusUnprocessableEntity, "bankdetails.tmpl", data)
		return
	}
	
	err=app.faculty.InsertBankDetails(id,form.BankName,form.AccountNumber,form.IFSC,"Insert picture")
	if err!=nil{
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "You have added your bank details successfully")
	http.Redirect(w, r, "/",http.StatusSeeOther)
}
func (app *application) facultyLoginPost(w http.ResponseWriter, r *http.Request) {
	var form facultyLoginForm

	err:=app.decodePostForm(r, &form)
	if err!=nil{
		app.clientError(w, http.StatusBadRequest)
		return
	}
        form.CheckField(validator.NotBlank(strconv.Itoa(form.FacultyID)), "facultyid", "This field cannot be blank")
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
			form.AddNonFieldError("ID or password is incorrect")
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

func (app *application) qpkCreate(w http.ResponseWriter, r *http.Request) {
        data:=app.newTemplateData(r)
        data.Form = qpkCreateForm{}
        app.render(w, http.StatusOK, "qpk.tmpl", data)
}

func (app *application) qpkCreatePost (w http.ResponseWriter, r *http.Request) {
	id:=app.sessionManager.Get(r.Context(),"authenticatedFacultyID").(int)
	var form qpkCreateForm

	err:=app.decodePostForm(r, &form)

        form.CheckField(validator.NotBlank(form.CourseCode), "coursecode", "This field must not be blank")
        form.CheckField(validator.NotBlank(strconv.Itoa(form.QuestionPaperCount)), "qc", "This field must be a valid number")
        form.CheckField(validator.NotBlank(strconv.Itoa(form.KeyCount)), "kc", "This field must be a valid number")

	if !form.Valid(){
		data:=app.newTemplateData(r)
		data.Form=form
		app.render(w, http.StatusUnprocessableEntity, "qpk.tmpl", data)
		return
	}
	
	tid,err:=app.honorarium.InsertQPK(id, form.CourseCode,form.QuestionPaperCount, form.KeyCount)
	if err!=nil{
		app.serverError(w, err)
		return
	}
	fmt.Println(tid)
	app.sessionManager.Put(r.Context(), "flash", "You have created a honorarium successfully")
	http.Redirect(w, r, "/",http.StatusSeeOther)
}


func (app *application) ansvCreate(w http.ResponseWriter, r *http.Request) {
        data:=app.newTemplateData(r)
        data.Form = ansvCreateForm{}
        app.render(w, http.StatusOK, "ansv.tmpl", data)
}

func (app *application) ansvCreatePost (w http.ResponseWriter, r *http.Request) {
	id:=app.sessionManager.Get(r.Context(),"authenticatedFacultyID").(int)
	var form ansvCreateForm

	err:=app.decodePostForm(r, &form)

        form.CheckField(validator.NotBlank(form.CourseCode), "coursecode", "This field must not be blank")
        form.CheckField(validator.NotBlank(strconv.Itoa(form.AnswerScriptCount)), "kc", "This field must be a valid number")

	if !form.Valid(){
		data:=app.newTemplateData(r)
		data.Form=form
		app.render(w, http.StatusUnprocessableEntity, "ansv.tmpl", data)
		return
	}
	
	tid,err:=app.honorarium.InsertValuedPaper(id, form.CourseCode,form.AnswerScriptCount)
	if err!=nil{
		app.serverError(w, err)
		return
	}
	fmt.Println(tid)
	app.sessionManager.Put(r.Context(), "flash", "You have created a honorarium successfully")
	http.Redirect(w, r, "/",http.StatusSeeOther)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
/* ADMIN HANDLERS */
