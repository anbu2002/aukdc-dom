package models
import (
	"database/sql"
	"errors"
)
type OtherModel struct {
	DB *sql.DB
} 
type Department struct {
	Degree string
	Branch string
	DegreeType string
	Department string
}

type Course struct {
	Department
	CourseCode string
	Title string
	Regulation string
}



func (m *OtherModel) GetCourse(coursecode string) (*Course, error) {
	s := &Course{
		Department: Department{
		},
	}

        err:= m.DB.QueryRow(`SELECT "Degree","Branch","DegreeType","Department",Course."CourseCode","Title","Regulation" FROM Department FULL JOIN Course ON Course."CourseCode"=department."CourseCode" WHERE (Course."CourseCode"=$1);`,coursecode).Scan(&s.Degree,&s.Branch,&s.DegreeType,&s.Department.Department,&s.CourseCode,&s.Title,&s.Regulation)
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

