package errare

import "fmt"

type ResponseError struct {
	status string
	client string
}

func NewResponseError(status string, clientName string) ResponseError {
	return ResponseError{
		status: status,
		client: clientName,
	}
}

func (e ResponseError) Error() string {
	client := "Server"
	if e.client != "" {
		client = e.client
	}
	return fmt.Sprintf("%s responded with error: %s", client, e.status)
}

type RequestError struct {
	err    error
	client string
}

func NewRequestError(err error, clientName string) RequestError {
	return RequestError{
		err:    err,
		client: clientName,
	}
}

func (e RequestError) Error() string {
	return fmt.Sprintf("Making request to %s failed: %s", e.client, e.err.Error())
}
