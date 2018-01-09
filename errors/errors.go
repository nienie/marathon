package errors

import (
	"fmt"
	"regexp"
)

var (
	errorConverters = make([]ErrorConverter, 0)
)

//ErrorConverter ...
type ErrorConverter func(error)ClientError


var defaultErrorConverter = func (err error) ClientError {
	if err == nil {
		return NewClientError(OK, nil)
	}
	str := err.Error()

	re := regexp.MustCompile(`getsockopt: connection refused`)
	if re.MatchString(str) {
		return NewClientError(ConnectException, err)
	}

	re = regexp.MustCompile(`dial.*i/o timeout`)
	if re.MatchString(str) {
		return NewClientError(SocketTimeoutException, err)
	}

	re = regexp.MustCompile(`read.*i/o timeout`)
	if re.MatchString(str) {
		return NewClientError(ReadTimeoutException, err)
	}

	return NewClientError(General, err)
}

//ClientError ...
type ClientError struct {
	errorType ErrorType
	err       error
}

//NewClientError ...
func NewClientError(errorType ErrorType, err error) ClientError {
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

//RegisterErrorConverters ...
func RegisterErrorConverters(converters ...ErrorConverter) {
	errorConverters = append(errorConverters, converters...)
}

//ConvertError ...
func ConvertError(err error) ClientError {
	for _, converter := range errorConverters {
		clientErr := converter(err)
		if clientErr.GetErrType() != OK {
			return clientErr
		}
	}

	return defaultErrorConverter(err)
}
