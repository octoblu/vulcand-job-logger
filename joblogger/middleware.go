package joblogger

import (
	"fmt"
	"net/http"
)

// Middleware is a vulcand middleware that logs to redis
type Middleware struct {
	RedisURI, RedisQueueName, BackendID string
}

// NewMiddleware constructs new Middleware instances
func NewMiddleware(RedisURI, RedisQueueName, BackendID string) (*Middleware, error) {
	if RedisURI == "" || RedisQueueName == "" || BackendID == "" {
		return nil, fmt.Errorf("RedisURI, RedisQueueName, and BackendID are all required. received '%v', '%v', and '%v'", RedisURI, RedisQueueName, BackendID)
	}

	return &Middleware{RedisURI, RedisQueueName, BackendID}, nil
}

// NewHandler returns a new Handler instance
func (middleware *Middleware) NewHandler(next http.Handler) (http.Handler, error) {
	return NewHandler(middleware.RedisURI, middleware.RedisQueueName, middleware.BackendID, next), nil
}

// String will be called by loggers inside Vulcand and command line tool.
func (middleware *Middleware) String() string {
	return fmt.Sprintf("redis-uri=%v, redis-queue-name=%v", middleware.RedisURI, middleware.RedisQueueName)
}
