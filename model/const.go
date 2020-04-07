package model

// Request Header keys
const (
	HeaderKeyUserAgent  = "User-Agent"
	HeaderKeyUserIP     = "X-Forwarded-For"
	HeaderKeyRequestID  = "requestid"
)

// Logger additional fields key
const (
	LoggerKeyRequestID  = "REQUEST_ID"
	LoggerKeyOperation  = "OPERATION"
	LoggerKeyUserIP     = "USER_IP"
	LoggerKeyUserAgent  = "USER_AGENT"
	ContextLogger       = "contextLogger"
	ContextHeader       = "contextHeader"
)
