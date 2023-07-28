package models
import (
	"database/sql"
	"errors"
)
type OtherModel struct {
	DB *sql.DB
}
type OfferedIn struct{
	DepartmentName string
}
type Programme struct {
	OfferedIn
	Degree string
	Branch string
	DegreeType string
}

type Course struct {
	Programme
	CourseCode string
	Title string
	Regulation string
}



func (m *OtherModel) GetCourse(coursecode string) (*Course, error) {
	s := &Course{}

	err:= m.DB.QueryRow(`SELECT "Degree","Branch","DegreeType","DepartmentName","CourseCode","Title","Regulation" FROM co_offeredin_pro WHERE "CourseCode"=$1;`,coursecode).Scan(&s.Degree,&s.Branch,&s.DegreeType,&s.DepartmentName,&s.CourseCode,&s.Title,&s.Regulation)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}
func (m *OtherModel) GetProgramme(branch string) (*Programme, error){
	s := &Programme{}
	
	err:= m.DB.QueryRow(`SELECT * FROM Programme WHERE "Branch" = $1`, branch).Scan(&s.Degree, &s.Branch, &s.DepartmentName, &s.DegreeType)
	if err !=nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *OtherModel) GetAllProgrammes() ([]*Programme, error){
	programmes:=[]*Programme{}
	rows, err:=m.DB.Query(`SELECT "Degree","Branch","DegreeType" FROM Programme`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next(){
		s:=&Programme{}
		err=rows.Scan(&s.Degree, &s.Branch, &s.DegreeType)
		if err != nil {
			return nil, err
		}
		programmes=append(programmes,s)
	}
	if err=rows.Err();err!=nil{
		return nil, err
	}
	return programmes, nil
}
func (m *OtherModel) GetAllCourseCodes() ([]*Course, error){
	courses:=[]*Course{}
	rows, err:=m.DB.Query(`SELECT "CourseCode" FROM Course`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next(){
		s:=&Course{}
		err=rows.Scan(&s.CourseCode)
		if err != nil {
			return nil, err
		}
		courses=append(courses,s)
	}
	if err=rows.Err();err!=nil{
		return nil, err
	}
	return courses, nil
}

