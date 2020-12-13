package errors

type ErrorCode string

const (
	InvalidRegisterData    ErrorCode = "400001"
	InvalidLoginData       ErrorCode = "400002"
	InvalidLoginCredential ErrorCode = "401001"
	UserNotFound           ErrorCode = "404001"
	RefreshTokenNotFound   ErrorCode = "404002"
	BearerTokenNotFound    ErrorCode = "404003"
	UserSignUpDataInvalid  ErrorCode = "422001"
	DatabaseQueryFailed    ErrorCode = "500001"
	BcryptProccessFaild    ErrorCode = "500002"
	UserLoginFailed        ErrorCode = "500003"
	TokenRefreshFailed     ErrorCode = "500004"
)
