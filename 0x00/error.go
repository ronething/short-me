package main

type Error interface {
	error
	Status() int
}

type StatusError struct {
	Code int
	Err error
}

func (se StatusError) Error() string {
	return se.Err.Error()
}

func (se StatusError) Stauts() int  {
	return se.Code
}

