package errs

import (
	"fmt"
)

func FailOnErr(err error, msg string) {
	if err != nil {
		panic(fmt.Errorf("%s %v", msg, err))
	}
}

type CredentialsError struct {
	Handle string
}

func (e *CredentialsError) Error() string {
	return fmt.Sprintf("Incorrect credentials for: %s", e.Handle)
}

func NewCredentialsErr(handle string) *CredentialsError {
	return &CredentialsError{
		Handle: handle,
	}
}
