package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"aukdc.dom.com/internal/models"
	"aukdc.dom.com/internal/validator"

	"github.com/julienschmidt/httprouter"
)

// Please change the data types of picture variables.
type facultySignupForm struct {
	FacultyID           int    `form:"facultyid"`
	Name                string `form:"name"`
	Phone               int64  `form:"phone"`
	Email               string `form:"email"`
	FacultyType         string `form:"facultytype"`
	Department          string `form:"dept"`
	Designation         string `form:"designation"`
	Password            string `form:"password"`
	PanID               string `form:"panid"`
	PanPicture          any    `form:"panpic"`
	Extension           int64  `form:"extnumber"`
	Esign               any    `form:"esign"`
	validator.Validator `form:"-"`
}

type bankDetailsForm struct {
	BankName            string `form:"bankname"`
	AccountNumber       int    `form:"accountno"`
	IFSC                string `form:"IFSC"`
	Passbook            any    `form:"passbook"`
	validator.Validator `form:"-"`
}
type userLoginForm struct {
	UserID              int    `form:"userid"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type honorariumCreateForm struct {
	TransactionID       int       `form:"-"`
	FacultyID           int       `form:"-"`
	CourseCode          string    `form:"coursecode"`
	InitialAmount       int       `form:"-"`
	FinalAmount         int       `form:"-"`
	CreatedTime         time.Time `form:"-"`
	validator.Validator `form:"-"`
}
type qpkCreateForm struct {
	honorariumCreateForm
	QuestionPaperCount int     `form:"qc"`
	KeyCount           int     `form:"kc"`
	QuestionPaperRate  float32 `form:"-"`
	KeyRate            float32 `form:"-"`
}
type ansvCreateForm struct {
	honorariumCreateForm
	AnswerScriptCount int     `form:"ac"`
	AnswerScriptRate  float32 `form:"-"`
}

/* COMMON HANDLERS */
//Method Stub
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "home.tmpl", data)
}
func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	if app.isAuthenticated(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.tmpl", data)
}
func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(validator.NotBlank(strconv.Itoa(form.UserID)), "userid", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}
	id, err := app.user.Authenticate(form.UserID, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("User ID or password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		} else {
			app.serverError(w, err)
		}
		return
	}
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	authority, err := app.user.Authorized(id)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	if authority {
		app.sessionManager.Put(r.Context(), "authorizedUserID", id)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Remove(r.Context(), "authorizedUserID")
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) facultyView(w http.ResponseWriter, r *http.Request) {
	id := app.sessionManager.Get(r.Context(), "authenticatedUserID").(int)
	var err error
	if app.isAuthorized(r) {
		params := httprouter.ParamsFromContext(r.Context())
		if params == nil {
			id = 0
		} else {
			id, err = strconv.Atoi(params.ByName("fid"))
			if err != nil {
				app.notFound(w)
				return
			}
		}
	}
	data := app.newTemplateData(r)
	if id != 0 {
		faculty, err := app.user.GetFaculty(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.notFound(w)
			} else {
				app.serverError(w, err)
			}
			return
		}
		data.Faculty = faculty
	}
	honoraria, err := app.honorarium.ViewAll(id)
	if err != nil {
		app.serverError(w, err)
		return
	}
	data.Honoraria = honoraria
	app.render(w, http.StatusOK, "honoraria.tmpl", data)
}
func (app *application) honorariumView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	data := app.newTemplateData(r)

	hid := params.ByName("hid")
	fid := app.sessionManager.Get(r.Context(), "authenticatedUserID").(int)
	var tyid int
	var err error
	if !app.isAuthorized(r){
		tyid,err =strconv.Atoi(params.ByName("tyid"))
		if err != nil {
			app.serverError(w, err)
			return
		}
	}else{
		honorarium, err:= app.honorarium.GetTransactionAdmin(hid)
		if err != nil {
			app.serverError(w, err)
			return
		}
		tyid=honorarium.TypeID
		fid=honorarium.FacultyID
	} 
	if tyid==1{
		qpk, err := app.honorarium.GetQPK(fid, hid)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.notFound(w)
				return
			} else {
			app.serverError(w, err)
			return
			}
		}
		data.QPK=qpk
		data.Honorarium=&qpk.Honorarium
	}else{
		vp, err := app.honorarium.GetValuedPaper(fid, hid)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.notFound(w)
				return
			} else {
			app.serverError(w, err)
			return
			}
		}
		data.VP=vp
		data.Honorarium=&vp.Honorarium
	}


	app.render(w, http.StatusOK, "honorarium.tmpl", data)
}
func (app *application) generatePrint(w http.ResponseWriter, r *http.Request){
	data := app.newTemplateData(r)
	fid:=app.sessionManager.Get(r.Context(), "authenticatedUserID").(int)
	var err error
	if app.isAuthorized(r){
		params:= httprouter.ParamsFromContext(r.Context())
		fid,err=strconv.Atoi(params.ByName("fid"))
		if err!=nil{
			app.serverError(w,err)
			return
		}
	}

	faculty,err :=app.user.GetFaculty(fid)
	if err!=nil{
		app.serverError(w,err)
		return
	}
	data.Faculty=faculty
	
	var honorarium *models.Honorarium
	params := httprouter.ParamsFromContext(r.Context())
	hid := params.ByName("hid")
	honorarium, err = app.honorarium.GetTransaction(fid, hid)
	if err!=nil{
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
			return
		} else {
		app.serverError(w, err)
		return
		}
	}
	data.Honorarium=honorarium

	course,err:=app.other.GetCourse(honorarium.CourseCode)
	if err!=nil{
		app.serverError(w,err)
		return
	}
	data.Course=course

	if honorarium.TypeID==1{
		qpk,err :=app.honorarium.GetQPK(fid,hid)
		if err!=nil{
			app.serverError(w,err)
			return
		}
		data.QPK=qpk
	}else{
		vp,err :=app.honorarium.GetValuedPaper(fid,hid)
		if err!=nil{
			app.serverError(w,err)
			return
		}
		data.VP=vp
	}

	bDetails,err:=app.user.GetBankDetails(fid)
	if err!=nil{
		app.serverError(w,err)
		return
	}
	data.BankDetails=bDetails

	app.render(w, http.StatusOK, "print-"+strconv.Itoa(honorarium.TypeID)+".tmpl", data)
}
/* FACULTY HANDLERS */
func (app *application) facultySignup(w http.ResponseWriter, r *http.Request) {
	if app.isAuthenticated(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	data := app.newTemplateData(r)
	data.Form = facultySignupForm{}
	app.render(w, http.StatusOK, "signup.tmpl", data)

}

func (app *application) facultySignupPost(w http.ResponseWriter, r *http.Request) {
	var form facultySignupForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
//must be in format AAAPZ1234C
	form.CheckField(validator.Matches(form.PanID, validator.PanRX), "panid", "This field must be a valid PAN ID")
	form.CheckField(validator.NotBlank(strconv.Itoa(form.FacultyID)), "facultyid", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Department), "dept", "This field cannot be blank")
	form.CheckField(validator.IntegerRange(form.Phone, 6000000000,9999999999), "phone", "This field must be a valid phone number")
	form.CheckField(validator.IntegerRange(form.Extension, 20000000,29999999), "extnumber", "This field must be a valid extension number")
	form.CheckField(validator.NotBlank(form.Designation), "designation", "Please select an appropriate designation")
	form.CheckField(validator.NotBlank(form.FacultyType), "facultytype", "Please select an appropriate type")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	PanPicture,err := app.uploadImage(w, r,"panpic",form.FacultyID)
	if err != nil {
		fmt.Println("handler", form.PanPicture)
		app.serverError(w, err)
		return
	}
	Esign,err := app.uploadImage(w, r,"esign",form.FacultyID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.user.Insert(form.FacultyID, form.Name, form.Phone, form.Email, form.FacultyType, form.Department, form.Designation, form.Password, form.PanID, PanPicture, form.Extension, Esign)
	if err != nil {
//add error for duplicate id
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "You have sucessfully signed up. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
func (app *application) addBankDetails(w http.ResponseWriter, r *http.Request) {
	if app.hasBankDetails(r) {
		app.sessionManager.Put(r.Context(), "flash", "Please request your admin to change your bank details")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	data := app.newTemplateData(r)
	data.Form = bankDetailsForm{}
	app.render(w, http.StatusOK, "bankdetails.tmpl", data)
}

func (app *application) addBankDetailsPost(w http.ResponseWriter, r *http.Request) {
	id := app.sessionManager.Get(r.Context(), "authenticatedUserID").(int)
	var form bankDetailsForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(validator.NotBlank(form.BankName), "bank", "This field must not be blank")
	form.CheckField(validator.NotBlank(strconv.Itoa(form.AccountNumber)), "accountno", "This field must be a valid account number")
//must be in format ABCD0678901
	form.CheckField(validator.Matches(form.IFSC, validator.IFSCRX), "IFSC", "This field must be a valid IFSC code")
	fmt.Println("handler", form.Passbook)

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "bankdetails.tmpl", data)
		return
	}

	Passbook,err := app.uploadImage(w, r,"passbook", id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.user.InsertBankDetails(id, form.BankName, form.AccountNumber, form.IFSC, Passbook)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "You have added your bank details successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) qpkCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	courses, err := app.other.GetAllCourses()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data.Courses = courses
	data.Form = qpkCreateForm{}
	app.render(w, http.StatusOK, "qpk.tmpl", data)
}

func (app *application) qpkCreatePost(w http.ResponseWriter, r *http.Request) {
	id := app.sessionManager.Get(r.Context(), "authenticatedUserID").(int)
	var form qpkCreateForm

	err := app.decodePostForm(r, &form)
	form.CheckField(validator.NotBlank(form.CourseCode), "coursecode", "This field must not be blank")
	form.CheckField(validator.NotBlank(strconv.Itoa(form.QuestionPaperCount)), "qc", "This field must be a valid number")
	form.CheckField(validator.NotBlank(strconv.Itoa(form.KeyCount)), "kc", "This field must be a valid number")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "qpk.tmpl", data)
		return
	}
	
	faculty,err:=app.user.GetFaculty(id)
	if err !=nil{
		app.serverError(w,err)
		return
	}

	tid, err := app.honorarium.InsertQPK(id, form.CourseCode, form.QuestionPaperCount, form.KeyCount, faculty.TDS)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "You have created a honorarium successfully")
	http.Redirect(w, r, "/honorarium/view/1/"+tid, http.StatusSeeOther)
}

func (app *application) ansvCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	courses, err := app.other.GetAllCourses()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data.Courses = courses
	data.Form = qpkCreateForm{}
	app.render(w, http.StatusOK, "ansv.tmpl", data)
}

func (app *application) ansvCreatePost(w http.ResponseWriter, r *http.Request) {
	id := app.sessionManager.Get(r.Context(), "authenticatedUserID").(int)
	var form ansvCreateForm

	err := app.decodePostForm(r, &form)

	form.CheckField(validator.NotBlank(form.CourseCode), "coursecode", "This field must not be blank")
	form.CheckField(validator.NotBlank(strconv.Itoa(form.AnswerScriptCount)), "ac", "This field must be a valid number")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "ansv.tmpl", data)
		return
	}

	faculty,err:=app.user.GetFaculty(id)
	if err !=nil{
		app.serverError(w,err)
		return
	}


	tid, err := app.honorarium.InsertValuedPaper(id, form.CourseCode, form.AnswerScriptCount, faculty.TDS)
	fmt.Println(tid)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "You have created a honorarium successfully")
	http.Redirect(w, r, "/honorarium/view/2/"+tid, http.StatusSeeOther)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

/* ADMIN HANDLERS */
func (app *application) facultyViewAll(w http.ResponseWriter, r *http.Request) {
	faculties, err := app.user.ViewAllFaculty()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.Faculties = faculties

	app.render(w, http.StatusOK, "faculty.tmpl", data)
}
