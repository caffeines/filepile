package errors

type ErrorCode string

const (
	InvalidRegisterData    ErrorCode = "400001"
	InvalidLoginData       ErrorCode = "400002"
	InvalidLoginCredential ErrorCode = "401001"
	UserNotFound           ErrorCode = "404001"
	DatabaseQueryFailed    ErrorCode = "500001"
	BcryptProccessFaild    ErrorCode = "500002"
	UserLoginFailed        ErrorCode = "500003"
)
