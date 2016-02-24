package joblogger

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/octoblu/vulcand-job-logger/wrapper"
)

var redisManager struct {
	Once        sync.Once
	Mutex       sync.Mutex
	Connections map[string]redis.Conn
}

// Handler implements http.Handler
type Handler struct {
	redisURI       string
	redisQueueName string
	backendID      string
	next           http.Handler
}

// NewHandler constructs a new handler
func NewHandler(redisURI, redisQueueName, backendID string, next http.Handler) *Handler {
	return &Handler{redisURI, redisQueueName, backendID, next}
}

// ServeHTTP will be called each time the request
// hits the location with this middleware activated
func (handler *Handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	redisChannel := make(chan []byte, 1000)

	go handler.logRequest(redisChannel)
	wrapped := wrapper.New(rw, redisChannel, time.Now(), handler.backendID)
	handler.next.ServeHTTP(wrapped, r)
}

func (handler *Handler) logRequest(logChannel chan []byte) {
	redisConn, err := handler.redisConn()
	if err != nil {
		logError("handler.redisConn() failed: %v\n", err)
		return
	}
	logEntryBytes := <-logChannel
	_, err = redisConn.Do("LPUSH", handler.redisQueueName, logEntryBytes)
	logError("Redis LPUSH failed: %v\n", err)
}

func (handler *Handler) redisConn() (redis.Conn, error) {
	redisManager.Once.Do(func() {
		redisManager.Connections = make(map[string]redis.Conn)
	})

	redisManager.Mutex.Lock()
	key := fmt.Sprintf("%v/%v", handler.redisURI, handler.redisQueueName)
	conn, ok := redisManager.Connections[key]
	if ok {
		redisManager.Mutex.Unlock()
		return conn, nil
	}

	conn, err := redis.DialURL(handler.redisURI)
	if err != nil {
		redisManager.Mutex.Unlock()
		return nil, err
	}

	redisManager.Connections[key] = conn
	redisManager.Mutex.Unlock()
	return conn, nil
}

func logError(fmtMessage string, err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, fmtMessage, err.Error())
}
