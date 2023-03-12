package models
import (
	"database/sql"
	"errors"
	"time"
)

type Honorarium struct {
	TransactionID int
	FacultyID string
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

func (m *HonorariumModel) InsertQPK(transactionID int, facultyID,courseCode string, initialAmount,finalAmount,questionPaperCount,keyCount int,questionPaperRate,keyRate float32) (int, error) {
	var id int
	err:= m.DB.QueryRow(`With new_qpk as(INSERT INTO Honorarium("TransactionID", "FacultyID", "CourseCode", "InitialAmount", "FinalAmount", "TypeID", "CreatedTime") VALUES ($1, $2, $3, $4, $5, 1 ,NOW()::timestamp(0)) RETURNING Honorarium."TransactionID") INSERT INTO "Question Paper/Key"("TransactionID","TypeID","QuestionPaperCount","KeyCount", "KeyRate", "QuestionPaperRate") VALUES((SELECT "TransactionID" FROM new_qpk), 1,$6,$7,$8,$9) returning 'TransactionID';`, transactionID, facultyID, courseCode, initialAmount,finalAmount, questionPaperCount, keyCount, questionPaperRate, keyRate).Scan(&id)
	if err!=nil{
		return  0,err
	}

	return id, nil
}

func (m *HonorariumModel) InsertValuedPaper(transactionID int, facultyID,courseCode string, initialAmount,finalAmount, answerScriptCount int,answerScriptRate float32) (int, error) {
	var id int
	err:= m.DB.QueryRow(`With new_ap as(INSERT INTO Honorarium("TransactionID", "FacultyID", "CourseCode", "InitialAmount", "FinalAmount", "TypeID", "CreatedTime") VALUES ($1, $2, $3, $4, $5, 1 ,NOW()::timestamp(0)) RETURNING Honorarium."TransactionID") INSERT INTO "Paper Valuation"("TransactionID","TypeID","AnswerScriptCount","AnswerScriptRate") VALUES((SELECT "TransactionID" FROM new_qpk), 2 ,$6,$7) returning 'TransactionID';`, transactionID, facultyID, courseCode, initialAmount,finalAmount, answerScriptCount, answerScriptRate).Scan(&id)
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
func (m *HonorariumModel) Latest() ([]*Honorarium, error) {
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
