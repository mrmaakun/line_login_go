package main

import ()

type APIError struct {
	ErrorCode        string
	ErrorDescription string
}

func (e APIError) Error() string {
	errorString := e.ErrorCode + ": " + e.ErrorDescription
	return errorString
}
