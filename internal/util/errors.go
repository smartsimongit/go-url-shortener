package util

import "errors"

var (
	ErrIncorrectPostURL     = errors.New("incorrect POST requestJSON URL")
	ErrIncorrectLongURL     = errors.New("you send incorrect LongURL")
	ErrIDParamIsMissing     = errors.New("id is missing in parameters")
	ErrIncorrectJSONRequest = errors.New("incorrect json requestJSON")
	ErrCreatedResponse      = errors.New("error created responseJson")
)
