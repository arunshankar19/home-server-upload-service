package domain

import "errors"

var (
	ErrUnauthorizedUser = errors.New("unauthorized user")

	ErrInvalidUserID = errors.New("user_id is invalid")
	
	ErrNotFileOwner = errors.New("the user is not the owner of target file")

	ErrUploadEventNotPresent = errors.New("upload event is not present")
)
