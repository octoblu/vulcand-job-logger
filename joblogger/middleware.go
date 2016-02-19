package logger

import (
	"fmt"
	"net/http"
)

// Middleware is a vulcand middleware that logs to redis
type Middleware struct {
	RedisURI, RedisQueueName string
}

// NewMiddleware constructs new Middleware instances
func NewMiddleware(RedisURI, RedisQueueName string) (*Middleware, error) {
	return &Middleware{RedisURI, RedisQueueName}, nil
}

// NewHandler returns a new Handler instance
func (middleware *Middleware) NewHandler(next http.Handler) (http.Handler, error) {
	return NewHandler(middleware.RedisURI, middleware.RedisQueueName, next), nil
}

// String will be called by loggers inside Vulcand and command line tool.
func (middleware *Middleware) String() string {
	return fmt.Sprintf("redis-uri=%v, redis-queue-name=%v", middleware.RedisURI, middleware.RedisQueueName)
}
