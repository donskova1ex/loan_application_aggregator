package internal

import "errors"

var (
	ErrRecordNoFound           = errors.New("no record found")
	ErrPhoneFormat             = errors.New("invalid phone format")
	ErrPhoneNumberExistToday   = errors.New("phone number already exists today")
	ErrInvalidPhoneNumber      = errors.New("invalid phone number")
	ErrEmptyPhoneNumber        = errors.New("empty phone number")
	ErrInvalidOrganizationName = errors.New("invalid organization name")
	ErrInvalidLoanApplication  = errors.New("invalid loan application")
)
