package http

const (
	MethodGet     = "GET"
	MethdodPost   = "POST"
	MethodPut     = "PUT"
	MethodDelete  = "DELETE"
	MethodHead    = "HEAD"
	MethodOptions = "OPTIONS"
)

var supportedMethods = map[string]struct{}{
	MethodGet:     {},
	MethdodPost:   {},
	MethodPut:     {},
	MethodDelete:  {},
	MethodHead:    {},
	MethodOptions: {},
}

func IsSupportedMethod(method string) bool {
	_, ok := supportedMethods[method]

	return ok
}
