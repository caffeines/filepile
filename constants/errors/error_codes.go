package errors

type ErrorCode string

const (
	InvalidRegisterData ErrorCode = "400001"
	DatabaseQueryFailed ErrorCode = "500001"
	BcryptProccessFaild ErrorCode = "500002"
)
