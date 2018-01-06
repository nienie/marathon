package errors

import "fmt"

//ClientError ...
type ClientError struct {
	errorType ErrorType
	err       error
}

//NewClientError ...
func NewClientError(errorType ErrorType, err error) error {
	return ClientError{
		errorType: errorType,
		err:       err,
	}

}

//Error ...
func (o ClientError) Error() string {
	if o.err != nil {
		return fmt.Sprintf("code = %d, type = %s, msg= %s", o.errorType, o.errorType.GetName(), o.err.Error())
	}
	return fmt.Sprintf("code = %d, type = %s, msg= %s", o.errorType, o.errorType.GetName(), "")
}

//GetErrType ...
func (o ClientError) GetErrType() ErrorType {
	return o.errorType
}
