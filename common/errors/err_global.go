package errors

// global errors
var (
	OK                  = NewError(0, "ok")
	ErrBadRequest       = NewError(1001001, "bad request")
	ErrServerError      = NewError(1001002, "server error")
	ErrServerBusy       = NewError(1001003, "server busy")
	ErrAuthInvalid      = NewError(1001004, "auth invalid")
	ErrAuthExpired      = NewError(1001005, "auth expired")
	ErrAuthForbiden     = NewError(1001006, "auth forbiden")
	ErrClientDeprecated = NewError(1001007, "this version of client is deprecated")
	ErrRateLimit        = NewError(1001008, "server busy")
)
