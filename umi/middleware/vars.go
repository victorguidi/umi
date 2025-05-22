package middleware

// This are not enabled by default, unless specified with:
// WithFlexibleCors function
var (
	FLEXIBLE_ORIGIN          string = "*"
	FLEXIBLE_COR_METHODS     string = "POST, GET, OPTIONS, PUT, DELETE, PATCH, TRACE, HEAD"
	FLEXIBLE_COR_HEADERS     string = "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"
	FLEXIBLE_COR_CREDENTIALS string = "true"
)
