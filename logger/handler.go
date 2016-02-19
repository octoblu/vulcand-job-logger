package logger

import "net/http"

// Handler implements http.Handler
type Handler struct {
	redisURI       string
	redisQueueName string
	next           http.Handler
}

// NewHandler constructs a new handler
func NewHandler(redisURI, redisQueueName string, next http.Handler) *Handler {
	return &Handler{redisURI, redisQueueName, next}
}

// ServeHTTP will be called each time the request
// hits the location with this middleware activated
func (handler *Handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	handler.next.ServeHTTP(rw, r)
}
