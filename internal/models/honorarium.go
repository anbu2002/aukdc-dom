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
	Branch string
	CourseCode string
	InitialAmount float32
	FinalAmount float32
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

func (m *HonorariumModel) InsertQPK(facultyID int,courseCode, branch string, questionPaperCount,keyCount int, tds float32) (string, error) {
	var initialAmount, finalAmount,questionPaperRate,keyRate float32
	questionPaperRate,keyRate=2000.0,3000.0
	initialAmount=questionPaperRate*float32(questionPaperCount)+keyRate*float32(keyCount)
	finalAmount=initialAmount-initialAmount*tds

	if finalAmount>5000{
		return "",ErrExceed
	}
	tid:=uuid.New()
	var id string
	err:= m.DB.QueryRow(`With new_qpk as(INSERT INTO honorarium("TransactionID", "FacultyID", "CourseCode", "InitialAmount", "FinalAmount", "TypeID", "CreatedTime","Branch") VALUES ($1, $2, $3, $4, $5, 1 ,NOW()::timestamp(0),$6) RETURNING honorarium."TransactionID") INSERT INTO "Question Paper/Key"("TransactionID","TypeID","QuestionPaperCount","KeyCount", "KeyRate", "QuestionPaperRate") VALUES((SELECT "TransactionID" FROM new_qpk), 1,$7,$8,$9,$10) returning "TransactionID";`, tid, facultyID, courseCode, initialAmount,finalAmount, branch, questionPaperCount, keyCount,keyRate, questionPaperRate).Scan(&id)
	if err!=nil{
		return  "", err
	}

	return id, nil
}

func (m *HonorariumModel) InsertValuedPaper(facultyID int, branch, courseCode string, answerScriptCount int, tds float32) (string, error) {
//ask criteria for change
	var answerScriptRate float32
	answerScriptRate=20.0
	initialAmount:=answerScriptRate*float32(answerScriptCount)
	finalAmount:=initialAmount-initialAmount*tds
	if(finalAmount<100){
		finalAmount=100
	}else if finalAmount>5000{
		return "",ErrExceed
	}
	tid:=uuid.New()
	var id string
	err:= m.DB.QueryRow(`With new_ap as(INSERT INTO honorarium("TransactionID", "FacultyID", "CourseCode", "InitialAmount", "FinalAmount", "TypeID", "CreatedTime","Branch") VALUES ($1, $2, $3, $4, $5, 2 ,NOW()::timestamp(0),$6) RETURNING honorarium."TransactionID") INSERT INTO "Paper Valuation"("TransactionID","TypeID","AnswerScriptCount","AnswerScriptRate") VALUES((SELECT "TransactionID" FROM new_ap), 2 ,$7,$8) returning "TransactionID";`, tid, facultyID, courseCode, initialAmount,finalAmount, branch, answerScriptCount, answerScriptRate).Scan(&id)
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
			TypeID: 1,
		},
	}

        err:= m.DB.QueryRow(`SELECT "CourseCode","InitialAmount","FinalAmount","QuestionPaperCount","KeyCount","KeyRate","QuestionPaperRate", "CreatedTime","Branch" FROM honorarium FULL JOIN "Question Paper/Key" ON honorarium."TransactionID"="Question Paper/Key"."TransactionID" WHERE (honorarium."TransactionID"=$1 AND "FacultyID"=$2)`,TransactionID, FacultyID).Scan(&s.CourseCode, &s.InitialAmount, &s.FinalAmount, &s.QuestionPaperCount, &s.KeyCount, &s.KeyRate, &s.QuestionPaperRate, &s.CreatedTime,&s.Branch)
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
			TypeID: 2,
		},
	}

        err:= m.DB.QueryRow(`SELECT "CourseCode","InitialAmount","FinalAmount","AnswerScriptRate","AnswerScriptCount", "CreatedTime","Branch" FROM honorarium FULL JOIN "Paper Valuation" ON honorarium."TransactionID"="Paper Valuation"."TransactionID" WHERE (honorarium."TransactionID"=$1 AND "FacultyID"=$2)`,TransactionID, FacultyID).Scan(&s.CourseCode, &s.InitialAmount, &s.FinalAmount, &s.AnswerScriptRate, &s.AnswerScriptCount, &s.CreatedTime,&s.Branch)
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

        err:= m.DB.QueryRow(`SELECT * FROM honorarium WHERE (honorarium."TransactionID"=$1 AND "FacultyID"=$2)`,TransactionID,FacultyID).Scan(&s.TransactionID,&s.FacultyID,&s.Branch,&s.CourseCode, &s.InitialAmount, &s.FinalAmount, &s.TypeID, &s.CreatedTime)
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

        err:= m.DB.QueryRow(`SELECT * FROM honorarium WHERE (honorarium."TransactionID"=$1)`,TransactionID).Scan(&s.TransactionID,&s.FacultyID,&s.Branch,&s.CourseCode, &s.InitialAmount, &s.FinalAmount, &s.TypeID, &s.CreatedTime)
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
		err=rows.Scan(&s.TransactionID,&s.FacultyID,&s.Branch,&s.CourseCode, &s.InitialAmount, &s.FinalAmount, &s.TypeID, &s.CreatedTime)
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
