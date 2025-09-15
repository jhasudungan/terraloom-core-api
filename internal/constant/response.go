package constant

const (

	// 01 : Resource Not Found
	// 02 : Authentication Failed
	// 03 : Access Denied
	// 04 : Validation failed
	// 05 : Conflict Resource
	// 99 : Unexpected Error

	SuccessCode    string = "00"
	SuccessMessage string = "Success"

	ResourceNotFoundCode    string = "01"
	ResourceNotFoundMessage string = "Resource Not Found"

	AuthFailedCode    string = "02"
	AuthFailedMessage string = "Authentication Failed"

	AccessDeniedCode    string = "03"
	AccessDeniedMessage string = "Access Denied"

	ValidationFailedCode    string = "04"
	ValidationFailedMessage string = "Validation failed"

	ConflictResourceCode    string = "05"
	ConflictResourceMessage string = "Conflict Resource"

	UnexpectedErrorCode    string = "99"
	UnexpectedErrorMessage string = "Unexpected Error"
)
