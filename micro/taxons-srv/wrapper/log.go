package wrapper

import (
	"log"
	"context"
	"github.com/micro/go-micro/server"
)

// LogWrapper is a log wrapper
func LogWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func (ctx context.Context, req server.Request, rsp interface{}) error {
		log.Printf("[Log Wrapper] Before serving request method: %v", req.Method())
		err := fn(ctx, req, rsp)
		log.Printf("[Log Wrapper] After serving request")
		return err
	}
}