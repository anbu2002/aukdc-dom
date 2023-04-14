package models
import (
	"database/sql"
	"errors"
	"time"

        "github.com/google/uuid"
)

type Honorarium struct {
	TransactionID int
	FacultyID int
	CourseCode string
	InitialAmount int
	FinalAmount int
	CreatedTime time.Time
}
type QPK struct {
	Honorarium
	QuestionPaperCount int
	KeyCount int
	QuestionPaperRate float32
	KeyRate float32
}
type ValuedPaper struct {
	Honorarium
	AnswerScriptCount int
	AnswerScriptRate float32
}
type HonorariumModel struct {
	DB *sql.DB
} 

func (m *HonorariumModel) InsertQPK(facultyID int,courseCode string, questionPaperCount,keyCount int) (int, error) {
	var initialAmount, finalAmount float64
	questionPaperRate,keyRate:=10.0, 100.0
	initialAmount=questionPaperRate*float64(questionPaperCount)+keyRate*float64(keyCount)
	finalAmount=initialAmount-initialAmount*0.10
	tid:=uuid.New()
	var id int
	err:= m.DB.QueryRow(`With new_qpk as(INSERT INTO honorarium("TransactionID", "FacultyID", "CourseCode", "InitialAmount", "FinalAmount", "TypeID", "CreatedTime") VALUES ($1, $2, $3, $4, $5, 1 ,NOW()::timestamp(0)) RETURNING honorarium."TransactionID") INSERT INTO "Question Paper/Key"("TransactionID","TypeID","QuestionPaperCount","KeyCount", "KeyRate", "QuestionPaperRate") VALUES((SELECT "TransactionID" FROM new_qpk), 1,$6,$7,$8,$9) returning 'TransactionID';`, tid, facultyID, courseCode, initialAmount,finalAmount, questionPaperCount, keyCount, questionPaperRate, keyRate).Scan(&id)
	if err!=nil{
		return  0,err
	}

	return id, nil
}

func (m *HonorariumModel) InsertValuedPaper(facultyID int,courseCode string, answerScriptCount int) (int, error) {
	answerScriptRate:=10.0
	initialAmount:=answerScriptRate*float64(answerScriptCount)
	finalAmount:=initialAmount-initialAmount*0.10

	tid:=uuid.New()
	var id int
	err:= m.DB.QueryRow(`With new_ap as(INSERT INTO honorarium("TransactionID", "FacultyID", "CourseCode", "InitialAmount", "FinalAmount", "TypeID", "CreatedTime") VALUES ($1, $2, $3, $4, $5, 1 ,NOW()::timestamp(0)) RETURNING honorarium."TransactionID") INSERT INTO "Paper Valuation"("TransactionID","TypeID","AnswerScriptCount","AnswerScriptRate") VALUES((SELECT "TransactionID" FROM new_qpk), 2 ,$6,$7) returning 'TransactionID';`, tid, facultyID, courseCode, initialAmount,finalAmount, answerScriptCount, answerScriptRate).Scan(&id)
	if err!=nil{
		return  0,err
	}

	return id, nil
}
/* Method Stub */
func (m *HonorariumModel) GetQPK(TransactionID int) (*QPK, error) {
	s := &QPK{}

	err:= m.DB.QueryRow(" ", TransactionID).Scan()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

/* Method Stub */
func (m *HonorariumModel) GetValuedPaper(TransactionID int) (*ValuedPaper, error) {
	s := &ValuedPaper{}

	err:= m.DB.QueryRow(" ", TransactionID).Scan()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

/* Method Stub */
func (m *HonorariumModel) ViewAll() ([]*Honorarium, error) {
	honoraria:= []*Honorarium{}
	rows, err:= m.DB.Query("")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next(){
		s:=&Honorarium{}
		err=rows.Scan()
	if err != nil {
		return nil, err
	}
	honoraria=append(honoraria,s)
	}
	if err=rows.Err();err!=nil{
		return nil, err
	}

	return honoraria, nil
}
