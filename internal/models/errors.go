package models

import(
	"errors"
)

var (
	ErrNoRecord=errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail = errors.New("models: duplicate email")
	ErrDuplicateID= errors.New("models: duplicate id")
	ErrDuplicatePhone= errors.New("models: duplicate phone")
	ErrDuplicateExtn= errors.New("models: duplicate extenstion")
	ErrDuplicatePan= errors.New("models: duplicate pan id")
	ErrDuplicateAccNo= errors.New("models: duplicate account number")
	ErrExceed=errors.New("models: exceeds 5000Rs.")
)


