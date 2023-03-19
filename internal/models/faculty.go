package models

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type Faculty struct{
	FacultyID string
	Name string
	Phone int64
	Email string
	FacultyType string
	Department string
	Designation string
	Password []byte
	PanID string
//	PanPicture
	Extension int64
//	Esign 
}

type FacultyModel struct{
	DB *sql.DB
}

func (m *FacultyModel) Insert(facultyID, name string, phoneNumber int64, email, facultyType, department, designation,  password, panID, panPicture string, extensionNumber int64, eSign string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO faculty ("FacultyID","Name", "PhoneNumber", "Email", "FacultyType", "Department", "Designation", "Password","PanID", "PanPicture", "ExtensionNumber", "Esign") VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_,err=m.DB.Exec(stmt, facultyID, name, phoneNumber, email, facultyType, department, designation, string(hashedPassword), panID, panPicture, extensionNumber, eSign)
//Add code for Faculty ID
	if err!=nil{
		var pSQLError *pq.Error
		if errors.As(err, &pSQLError){
			if pSQLError.Code == "23505" && strings.Contains(pSQLError.Message, "users_uc_email"){
			return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (m *FacultyModel) Authenticate(facultyID, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt:=`SELECT "FacultyID", "Password" FROM faculty where "FacultyID"=$1`

	err:=m.DB.QueryRow(stmt, facultyID).Scan(&id, &hashedPassword)
	if err!=nil{
		if errors.Is(err, sql.ErrNoRows){
			return 0, ErrInvalidCredentials
		}else {
			return 0, err
		}
	}
	
	err=bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err!=nil{
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword){
		return 0, ErrInvalidCredentials
	}else{
		return 0,err
		}
	}
	return id, nil
}

func (m *FacultyModel) Exists(id int) (bool, error) {
	var exists bool

	stmt:=`SELECT EXISTS(SELECT true FROM faculty WHERE "FacultyID"=$1)`

	err:=m.DB.QueryRow(stmt, id).Scan(&exists)

	return exists, err
}
