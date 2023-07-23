package main

import (
	"errors"
//	"fmt"
	"net/http"
	"strconv"
	"time"
	"strings"

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
	CPassword            string `form:"cpassword"`
	PanID               string `form:"panid"`
	PanPicture          any    `form:"panpic"`
	Extension           int64  `form:"extnumber"`
	Esign               any    `form:"esign"`
	bankDetailsForm
	validator.Validator `form:"-"`
}

type bankDetailsForm struct {
	BankName            string `form:"bankname"`
	AccountNumber       int64  `form:"accountno"`
	CAccountNumber       int64  `form:"caccountno"`
	IFSC                string `form:"IFSC"`
	Passbook            any    `form:"passbook"`
}
type userLoginForm struct {
	UserID              int    `form:"userid"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type honorariumCreateForm struct {
	TransactionID       int       `form:"-"`
	FacultyID           int       `form:"-"`
	Department          string    `form:"dept"`
	Branch              string    `form:"branch"`
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
	CourseCode2		   string  `form:"coursecode2"`
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
	if !app.isAuthorized(r) {
		tyid, err = strconv.Atoi(params.ByName("tyid"))
		if err != nil {
			app.serverError(w, err)
			return
		}
	} else {
		honorarium, err := app.honorarium.GetTransactionAdmin(hid)
		if err != nil {
			app.serverError(w, err)
			return
		}
		tyid = honorarium.TypeID
		fid = honorarium.FacultyID
	}
	if tyid == 1 {
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
		data.QPK = qpk
		data.Honorarium = &qpk.Honorarium
	} else {
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
		data.VP = vp
		data.Honorarium = &vp.Honorarium
	}
	data.Faculty, err = app.user.GetFaculty(fid)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, http.StatusOK, "honorarium.tmpl", data)
}
func (app *application) generatePrint(w http.ResponseWriter, r *http.Request) {
	//needs optimization
	data := app.newTemplateData(r)
	fid := app.sessionManager.Get(r.Context(), "authenticatedUserID").(int)
	var err error
	if app.isAuthorized(r) {
		params := httprouter.ParamsFromContext(r.Context())
		fid, err = strconv.Atoi(params.ByName("fid"))
		if err != nil {
			app.serverError(w, err)
			return
		}
	}

	faculty, err := app.user.GetFaculty(fid)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data.Faculty = faculty
	var honorarium *models.Honorarium
	params := httprouter.ParamsFromContext(r.Context())
	hid := params.ByName("hid")
	honorarium, err = app.honorarium.GetTransaction(fid, hid)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
			return
		} else {
			app.serverError(w, err)
			return
		}
	}
	data.Honorarium = honorarium

	courseStr := honorarium.CourseCode
	coursesStr := strings.Split(courseStr, ",")
	for _, courseCode := range coursesStr{
		course, err:=app.other.GetCourse(courseCode)
		if err!=nil{
			app.serverError(w,err)
			return
		}
		data.Courses = append(data.Courses, course)
	}

	programme,err := app.other.GetProgramme(data.Courses[0].Branch)
	if err!=nil{
		app.serverError(w,err)
		return
	}
	data.Programmes = append(data.Programmes, programme)

	if honorarium.TypeID == 1 {
		qpk, err := app.honorarium.GetQPK(fid, hid)
		if err != nil {
			app.serverError(w, err)
			return
		}
		data.QPK = qpk
	} else {
		vp, err := app.honorarium.GetValuedPaper(fid, hid)
		if err != nil {
			app.serverError(w, err)
			return
		}
		data.VP = vp
	}

	bDetails, err := app.user.GetBankDetails(fid)
	if err != nil {
		app.serverError(w, err)
		return
	}
	data.BankDetails = bDetails

	app.render(w, http.StatusOK, "print-"+strconv.Itoa(honorarium.TypeID)+".tmpl", data)
}

/* FACULTY HANDLERS */
func (app *application) addDetails(w http.ResponseWriter, r *http.Request) {
	if app.hasBankDetails(r) {
		app.sessionManager.Put(r.Context(), "flash", "Please request your admin to change your bank details")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	fid:=app.sessionManager.Get(r.Context(), "authenticatedUserID").(int)
	
	faculty, err := app.user.GetFacultyStage(fid)
	if err!=nil{
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Form = facultySignupForm{
			Name: faculty.Name,
			Designation: faculty.Designation,
			Department: faculty.DepartmentName,
		}
	app.render(w, http.StatusOK, "first-login.tmpl", data)
}

func (app *application) addDetailsPost(w http.ResponseWriter, r *http.Request) {
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
	form.CheckField(validator.IntegerRange(form.Phone, 6000000000, 9999999999), "phone", "This field must be a valid phone number")
	form.CheckField(validator.IntegerRange(form.Extension, 20000000, 29999999), "extnumber", "This field must be a valid extension")
	form.CheckField(validator.NotBlank(form.Designation), "designation", "Please select an appropriate designation")
	form.CheckField(validator.NotBlank(form.FacultyType), "facultytype", "Please select an appropriate type")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")
	form.CheckField(validator.PermittedValue(form.Password, form.CPassword), "password", "Password is not the same as confirm password")

	//to optmize
	form.CheckField(validator.NotBlank(form.BankName), "bankname", "This field must not be blank")
	form.CheckField(validator.MinChars(strconv.FormatInt(form.AccountNumber, 10), 10), "accountno", "This field must be a valid account number")
	form.CheckField(validator.MaxChars(strconv.FormatInt(form.AccountNumber, 10), 16), "accountno", "This field must be a valid account number")
	//must be in format SBIN0005943
	form.CheckField(validator.Matches(form.IFSC, validator.IFSCRX), "IFSC", "This field must be a valid IFSC code")
	form.CheckField(validator.PermittedValue(form.AccountNumber, form.CAccountNumber), "accountno", "Account Number does not match")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "first-login.tmpl", data)
		return
	}

	panPicture, _, err := r.FormFile("panpic")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			form.AddFieldError("panpic", "Please upload a valid picture")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "first-login.tmpl", data)
			return
		}
		app.serverError(w, err)
		return
	}
	defer panPicture.Close()

	esignPicture, _, err := r.FormFile("esign")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			form.AddFieldError("esign", "Please upload valid picture")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
			return
		}
		app.serverError(w, err)
		return
	}
	defer esignPicture.Close()
	pan,err := app.convertImage(panPicture)
	if err != nil{
		app.serverError(w, err)
		return
	}
	esign,err := app.convertImage(esignPicture)
	if err != nil{
		app.serverError(w, err)
		return
	}

	err = app.user.Insert(form.FacultyID, form.Name, form.Phone, form.Email, form.FacultyType, form.Department, form.Designation, form.Password, form.PanID, pan, form.Extension, esign)
	if err != nil {
		var flag bool
		if errors.Is(err, models.ErrDuplicateID) {
			form.AddFieldError("facultyid", "FacultyID is already in use")
			flag = true
		} else if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			flag = true
		} else if errors.Is(err, models.ErrDuplicatePhone) {
			form.AddFieldError("phone", "Phone Number is already in use")
			flag = true
		} else if errors.Is(err, models.ErrDuplicateExtn) {
			form.AddFieldError("extnumber", "Extension Number is already in use")
			flag = true
		} else if errors.Is(err, models.ErrDuplicatePan) {
			form.AddFieldError("panid", "PAN ID is already in use")
			flag = true
		} else if errors.Is(err, models.ErrInvalidDepartment) {
			form.AddFieldError("dept", "Please enter in the right department")
			flag = true
		}
		if flag {
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "first-login.tmpl", data)
			return
		}
		app.serverError(w, err)
		return
	}
	pb, _, err := r.FormFile("passbook")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			form.AddFieldError("passbook", "Please upload valid picture")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "first-login.tmpl", data)
			return
		}
		app.serverError(w, err)
		return
	}
	defer pb.Close()

	passbook,err := app.convertImage(pb)
	if err != nil{
		app.serverError(w, err)
		return
	}

	err = app.user.InsertBankDetails(form.FacultyID, form.BankName, form.AccountNumber, form.IFSC, passbook)
	if err != nil {
		app.serverError(w, err)
		return
	}
	_, err =app.user.RemoveFacultyStage(app.sessionManager.Get(r.Context(), "authenticatedUserID").(int))
	if err!=nil{
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "You have added your details successfully, please login again")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) qpkCreate(w http.ResponseWriter, r *http.Request) {
	//needs optimization
	data := app.newTemplateData(r)
	courses, err := app.other.GetAllCourseCodes()
	if err != nil {
		app.serverError(w, err)
		return
	}
	programmes, err := app.other.GetAllProgrammes()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data.Courses = courses
	data.Programmes = programmes
	data.Form = qpkCreateForm{}
	app.render(w, http.StatusOK, "qpk.tmpl", data)
}

func (app *application) qpkCreatePost(w http.ResponseWriter, r *http.Request) {
	id := app.sessionManager.Get(r.Context(), "authenticatedUserID").(int)
	var form qpkCreateForm

	err := app.decodePostForm(r, &form)
	courseCode := form.CourseCode
	form.CheckField(validator.NotBlank(courseCode), "coursecode", "This field must not be blank")
	form.CheckField(validator.NotBlank(form.Branch), "branch", "This field must not be blank")
	form.CheckField(validator.PermittedValue(form.QuestionPaperCount, 0, 1, 2), "qc", "This field must be a valid number")
	form.CheckField(validator.PermittedValue(form.KeyCount, 0, 1), "kc", "This field must be a valid number")
	if form.QuestionPaperCount == 2 {
		form.CheckField(validator.NotBlank(form.CourseCode2), "coursecode2", "This field must not be blank")
		courseCode = form.CourseCode + "," + form.CourseCode2
	}

	if !form.Valid() {
		app.sessionManager.Put(r.Context(), "flash", "Please enter in valid details")
		http.Redirect(w, r, "/honorarium/create/qpk", http.StatusSeeOther)
		return
	}

	faculty, err := app.user.GetFaculty(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	tid, err := app.honorarium.InsertQPK(id, courseCode, form.Branch, form.QuestionPaperCount, form.KeyCount, faculty.TDS)
	if err != nil {
		if errors.Is(err, models.ErrExceed) {
			app.sessionManager.Put(r.Context(), "flash", "Final Amount is 0 or exceeds Rs. 5000, please try again")
			http.Redirect(w, r, "/honorarium/create/qpk", http.StatusSeeOther)
			return
		}
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "You have created a honorarium successfully")
	http.Redirect(w, r, "/honorarium/view/1/"+tid, http.StatusSeeOther)
}

func (app *application) ansvCreate(w http.ResponseWriter, r *http.Request) {
	//needs optimization
	data := app.newTemplateData(r)
	courses, err := app.other.GetAllCourseCodes()
	if err != nil {
		app.serverError(w, err)
		return
	}
	programmes, err := app.other.GetAllProgrammes()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data.Courses = courses
	data.Programmes = programmes
	data.Form = ansvCreateForm{}
	app.render(w, http.StatusOK, "ansv.tmpl", data)
}

func (app *application) ansvCreatePost(w http.ResponseWriter, r *http.Request) {
	id := app.sessionManager.Get(r.Context(), "authenticatedUserID").(int)
	var form ansvCreateForm

	err := app.decodePostForm(r, &form)

	form.CheckField(validator.NotBlank(form.CourseCode), "coursecode", "This field must not be blank")
	form.CheckField(validator.NotBlank(form.Branch), "branch", "This field must not be blank")
	form.CheckField(validator.NotBlank(strconv.Itoa(form.AnswerScriptCount)), "ac", "This field must be a valid number")

	if !form.Valid() {
		app.sessionManager.Put(r.Context(), "flash", "Please enter in valid details")
		http.Redirect(w, r, "/honorarium/create/ansv", http.StatusSeeOther)
		return
	}

	faculty, err := app.user.GetFaculty(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	tid, err := app.honorarium.InsertValuedPaper(id, form.Branch, form.CourseCode, form.AnswerScriptCount, faculty.TDS)
	if err != nil {
		if errors.Is(err, models.ErrExceed) {
			app.sessionManager.Put(r.Context(), "flash", "Final Amount is 0 or exceeds Rs. 5000, please try again")
			http.Redirect(w, r, "/honorarium/create/ansv", http.StatusSeeOther)
			return
		}
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
