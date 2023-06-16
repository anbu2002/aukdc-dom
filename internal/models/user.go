package models

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)
type User struct{
	ID int
	Name string
	Phone int64
	Email string
	Password []byte
}
type BankDetails struct{
	BankName string
	AccountNumber int
	IFSC string
	Passbook string
}
type Faculty struct{
	User
	FacultyType string
	Department string
	Designation string
	PanID string
	PanPicture string
	Extension int64
	Esign string
	TDS float32
	BankDetails 
}

type UserModel struct{
	DB *sql.DB
}


func (m *UserModel) Insert(facultyID int, name string, phoneNumber int64, email, facultyType, department, designation,  password, panID, panPicture string, extensionNumber int64, eSign string) error {
	tdsper:=0.1
	if facultyType=="Visiting" || facultyType=="Contract/Guest" {
		tdsper=0.0
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `With new_user as (INSERT INTO Users("ID","Name","PhoneNumber","Email","HashedPassword","RoleID") VALUES($1,$2,$3,$4,$5,2) RETURNING users."ID") INSERT INTO faculty ("FacultyID","FacultyType", "Department", "Designation", "PanID", "PanPicture", "ExtensionNumber", "Esign", "TDS") VALUES((SELECT "ID" FROM new_User), $6, $7, $8, $9, $10, $11, $12,$13);`
	_,err=m.DB.Exec(stmt, facultyID, name, phoneNumber, email, string(hashedPassword),facultyType, department, designation, panID, panPicture, extensionNumber, eSign,tdsper)
	if err!=nil{
		var pSQLError *pq.Error
		if errors.As(err, &pSQLError){
			if pSQLError.Code == "23505"{
				if strings.Contains(pSQLError.Message, "users_Email_key"){
					return ErrDuplicateEmail
				}else if strings.Contains(pSQLError.Message, "users_pkey"){
					return ErrDuplicateID
				}else if strings.Contains(pSQLError.Message, "users_PhoneNumber_key"){
					return ErrDuplicatePhone
				}else if strings.Contains(pSQLError.Message, "faculty_ExtenstionNumber_key"){
					return ErrDuplicateExtn
				}else if strings.Contains(pSQLError.Message, "faculty_PanID_key"){
					return ErrDuplicatePan
				}
			}
		}
		return err
	}
	return nil
}

func (m *UserModel) InsertBankDetails(facultyID int, bankName string, accountNumber int, IFSC, passbook string) error{
	_,err:=m.DB.Exec(`INSERT INTO account("BankName","FacultyID","AccountNumber","IFSCCode","Passbook") VALUES ($1,$2,$3,$4,$5)`,bankName,facultyID,accountNumber,IFSC,passbook)
	if err!=nil{
		var pSQLError *pq.Error
		if errors.As(err, &pSQLError){
			if pSQLError.Code == "23505" && strings.Contains(pSQLError.Message, "account__pkey"){
				return ErrDuplicateAccNo
			}
		return err
		}
	}
	return nil
	
}


func (m *UserModel) GetBankDetails(facultyID int) (*BankDetails, error){
	b:=&BankDetails{}
	stmt:=`SELECT "BankName","AccountNumber","IFSCCode","Passbook" FROM Account WHERE "FacultyID"=$1`
	err:=m.DB.QueryRow(stmt,facultyID).Scan(&b.BankName, &b.AccountNumber, &b.IFSC, &b.Passbook)
	if err!=nil{
		if errors.Is(err,sql.ErrNoRows){
			return nil, ErrNoRecord
		} else{
			return nil, err
		}
	}
	return b, nil
}

func (m *UserModel) HasBankDetails(id int) (bool, error) {
	var exists bool

	stmt:=`SELECT EXISTS(SELECT true FROM account WHERE "FacultyID"=$1)`

	err:=m.DB.QueryRow(stmt, id).Scan(&exists)

	return exists, err
}
func (m *UserModel) Authenticate(facultyID int, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt:=`SELECT "ID", "HashedPassword" FROM users where "ID"=$1`

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

func (m *UserModel) Authorized(id int) (bool, error){
        var authority bool

	stmt:=`SELECT EXISTS(SELECT true FROM users WHERE ("ID"=$1 AND "RoleID"=1))`

        err:=m.DB.QueryRow(stmt, id).Scan(&authority)

        return authority, err
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool

	stmt:=`SELECT EXISTS(SELECT true FROM users WHERE "ID"=$1)`

	err:=m.DB.QueryRow(stmt, id).Scan(&exists)

	return exists, err
}


func (m *UserModel) ViewAllFaculty() ([]*Faculty, error) {
        faculties:= []*Faculty{}
        rows, err:= m.DB.Query(`SELECT "ID","Name","PhoneNumber","Email","FacultyType","Department","Designation","PanID","PanPicture","ExtensionNumber","Esign" ,"TDS" FROM users FULL JOIN Faculty ON "ID"="FacultyID" WHERE "RoleID"=2`)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        for rows.Next(){
                s:=&Faculty{}
                err=rows.Scan(&s.ID, &s.Name, &s.Phone, &s.Email, &s.FacultyType, &s.Department, &s.Designation, &s.PanID, &s.PanPicture, &s.Extension, &s.Esign,&s.TDS)
        if err != nil {
                return nil, err
        }
        faculties=append(faculties,s)
        }
        if err=rows.Err();err!=nil{
                return nil, err
        }

        return faculties, nil
}

func (m *UserModel) GetFaculty(fid int) (*Faculty, error) {
        s := &Faculty{}

        err:= m.DB.QueryRow(`SELECT "ID","Name","PhoneNumber","Email","FacultyType","Department","Designation","PanID","PanPicture","ExtensionNumber","Esign","TDS" FROM users FULL JOIN Faculty ON "ID"="FacultyID" WHERE "ID"=$1`,fid).Scan(&s.ID, &s.Name, &s.Phone, &s.Email, &s.FacultyType, &s.Department, &s.Designation, &s.PanID, &s.PanPicture, &s.Extension, &s.Esign,&s.TDS)
        if err != nil {
                if errors.Is(err, sql.ErrNoRows) {
                        return nil, ErrNoRecord
                } else {
                        return nil, err
                }
        }
        return s, nil
}
