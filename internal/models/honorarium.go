package models
import (
	"database/sql"
	"errors"
	"time"

        "github.com/google/uuid"
)

type Honorarium struct {
	TransactionID string
	FacultyID int
	CourseCode string
	InitialAmount int
	FinalAmount int
	TypeID int
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

func (m *HonorariumModel) InsertQPK(facultyID int,courseCode string, questionPaperCount,keyCount int) (string, error) {
	var initialAmount, finalAmount float64
	questionPaperRate,keyRate:=10.0, 100.0
	initialAmount=questionPaperRate*float64(questionPaperCount)+keyRate*float64(keyCount)
	finalAmount=initialAmount-initialAmount*0.10
	tid:=uuid.New()
	var id string
	err:= m.DB.QueryRow(`With new_qpk as(INSERT INTO honorarium("TransactionID", "FacultyID", "CourseCode", "InitialAmount", "FinalAmount", "TypeID", "CreatedTime") VALUES ($1, $2, $3, $4, $5, 1 ,NOW()::timestamp(0)) RETURNING honorarium."TransactionID") INSERT INTO "Question Paper/Key"("TransactionID","TypeID","QuestionPaperCount","KeyCount", "KeyRate", "QuestionPaperRate") VALUES((SELECT "TransactionID" FROM new_qpk), 1,$6,$7,$8,$9) returning "TransactionID";`, tid, facultyID, courseCode, initialAmount,finalAmount, questionPaperCount, keyCount, questionPaperRate, keyRate).Scan(&id)
	if err!=nil{
		return  "", err
	}

	return id, nil
}

func (m *HonorariumModel) InsertValuedPaper(facultyID int,courseCode string, answerScriptCount int) (string, error) {
	answerScriptRate:=10.0
	initialAmount:=answerScriptRate*float64(answerScriptCount)
	finalAmount:=initialAmount-initialAmount*0.10

	tid:=uuid.New()
	var id string
	err:= m.DB.QueryRow(`With new_ap as(INSERT INTO honorarium("TransactionID", "FacultyID", "CourseCode", "InitialAmount", "FinalAmount", "TypeID", "CreatedTime") VALUES ($1, $2, $3, $4, $5, 1 ,NOW()::timestamp(0)) RETURNING honorarium."TransactionID") INSERT INTO "Paper Valuation"("TransactionID","TypeID","AnswerScriptCount","AnswerScriptRate") VALUES((SELECT "TransactionID" FROM new_qpk), 2 ,$6,$7) returning "TransactionID";`, tid, facultyID, courseCode, initialAmount,finalAmount, answerScriptCount, answerScriptRate).Scan(&id)
	if err!=nil{
		return  "", err
	}

	return id, nil
}
func (m *HonorariumModel) GetQPK(FacultyID int ,TransactionID string) (*QPK, error) {
	s := &QPK{
		Honorarium: Honorarium{
			TransactionID: TransactionID,
			FacultyID: FacultyID,
		},
	}

        err:= m.DB.QueryRow(`SELECT "CourseCode","InitialAmount","FinalAmount","QuestionPaperCount","KeyCount","KeyRate","QuestionPaperRate", "CreatedTime" FROM honorarium FULL JOIN "Question Paper/Key" ON honorarium."TransactionID"="Question Paper/Key"."TransactionID" WHERE (honorarium."TransactionID"=$1 AND "FacultyID"=$2)`).Scan(&s.CourseCode, &s.InitialAmount, &s.FinalAmount, &s.QuestionPaperCount, &s.KeyCount, &s.KeyRate, &s.QuestionPaperRate, &s.CreatedTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *HonorariumModel) GetValuedPaper(FacultyID int, TransactionID string) (*ValuedPaper, error) {
	s := &ValuedPaper{
		Honorarium: Honorarium{
			TransactionID: TransactionID,
			FacultyID: FacultyID,
		},
	}

        err:= m.DB.QueryRow(`SELECT "CourseCode","InitialAmount","FinalAmount","AnswerScriptRate","AnswerScriptCount", "CreatedTime" FROM honorarium FULL JOIN "Question Paper/Key" ON honorarium."TransactionID"="Question Paper/Key"."TransactionID" WHERE (honorarium."TransactionID"=$1 AND "FacultyID"=$2)`).Scan(&s.CourseCode, &s.InitialAmount, &s.FinalAmount, &s.AnswerScriptRate, &s.AnswerScriptCount, &s.CreatedTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}
func (m *HonorariumModel) GetTransaction(FacultyID int ,TransactionID string) (*Honorarium, error) {
	s := &Honorarium{}

        err:= m.DB.QueryRow(`SELECT * FROM honorarium WHERE (honorarium."TransactionID"=$1 AND "FacultyID"=$2)`,TransactionID,FacultyID).Scan(&s.TransactionID,&s.FacultyID,&s.CourseCode, &s.InitialAmount, &s.FinalAmount, &s.TypeID, &s.CreatedTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}
func (m *HonorariumModel) GetTransactionAdmin(TransactionID string) (*Honorarium, error) {
	s := &Honorarium{}

        err:= m.DB.QueryRow(`SELECT * FROM honorarium WHERE (honorarium."TransactionID"=$1)`,TransactionID).Scan(&s.TransactionID,&s.FacultyID,&s.CourseCode, &s.InitialAmount, &s.FinalAmount, &s.TypeID, &s.CreatedTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}
func (m *HonorariumModel) ViewAll(FacultyID int) ([]*Honorarium, error) {
	honoraria:= []*Honorarium{}
	stmt:=`SELECT * FROM honorarium WHERE ("FacultyID"`
	if FacultyID==0{
		stmt+=`>$1)`
	}else{
		stmt+=`=$1)`
	}
	rows, err:= m.DB.Query(stmt,FacultyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next(){
		s:=&Honorarium{}
		err=rows.Scan(&s.TransactionID,&s.FacultyID,&s.CourseCode, &s.InitialAmount, &s.FinalAmount, &s.TypeID, &s.CreatedTime)
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
