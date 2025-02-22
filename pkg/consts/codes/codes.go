package codes

const (
	OK = 1000 + iota
	InternalError
	InvalidParams
	Unauthorized
	NotFound
	Forbidden
	AlreadyExists
	MethodNotAllowed
)
