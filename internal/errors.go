package internal

import "errors"

var (
	ErrRecordNoFound         = errors.New("no record found")
	ErrPhoneFormat           = errors.New("invalid phone format")
	ErrPhoneNumberExistToday = errors.New("phone number already exists today")
)
