package proctify

import "fmt"

type ProctifyError struct {
	Err     error
	Message string
}

func newError(e error) ProctifyError {
	return ProctifyError{
		Err:     e,
		Message: fmt.Sprintf("proctify error: %s", e.Error()),
	}
}

func (e ProctifyError) Error() string {
	return e.Message
}

func (e ProctifyError) Unwrap() error {
	return e.Err
}
