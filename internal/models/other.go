package models
import (
	"database/sql"
	"errors"
)
type OtherModel struct {
	DB *sql.DB
} 

type Course struct {
	Department
	CourseCode string
	Title string
	Regulation string
}

type Department struct {
	Degree string
	Branch string
	DegreeType string
	Department string
}


func (m *OtherModel) GetCourse(coursecode string) (*Course, error) {
	s := &Course{}

        err:= m.DB.QueryRow(`SELECT * FROM Course WHERE (Course."CourseCode"=$1)`,coursecode).Scan(&s.CourseCode,&s.Title,&s.Regulation)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *OtherModel) GetAllCourses() ([]*Course, error){
	courses:=[]*Course{}
	rows, err:=m.DB.Query(`SELECT * FROM course`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next(){
		s:=&Course{}
		err=rows.Scan(&s.CourseCode,&s.Title,&s.Regulation)
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

