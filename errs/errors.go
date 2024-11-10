package errs

import (
	"fmt"
)

func FailOnErr(err error, msg string) {
	if err != nil {
		panic(fmt.Errorf("%s %v", msg, err))
	}
}

type IncorrectCredentials struct {
	error
}

func (e IncorrectCredentials) Error() string {
	return "Incorrect Credentials Entered"
}

func NewIncorrectCredentialsError() IncorrectCredentials {
	return IncorrectCredentials{}
}
