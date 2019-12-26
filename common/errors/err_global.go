package errors

// error code format:
//       code = [appid(2)][module(3)][code(3)]
//
// global errors:
var (
	OK                  = NewError(0, "ok")
	ErrBadRequest       = NewError(10100001, "bad request")
	ErrServerError      = NewError(10100002, "server error")
	ErrServerBusy       = NewError(10100003, "server busy")
	ErrAuthInvalid      = NewError(10100004, "auth invalid")
	ErrAuthExpired      = NewError(10100005, "auth expired")
	ErrAuthForbiden     = NewError(10100006, "auth forbiden")
	ErrClientDeprecated = NewError(10100007, "this version of client is deprecated")
	ErrRateLimit        = NewError(10100008, "server busy")
)
