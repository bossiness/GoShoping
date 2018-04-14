package logwrapper

import (
	"github.com/micro/go-log"
	"context"
	"github.com/micro/go-micro/server"
)

// LogWrapper is a log wrapper
func LogWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func (ctx context.Context, req server.Request, rsp interface{}) error {
		log.Log("[Log Wrapper] Before serving request service: %v", req.Service())
		log.Log("[Log Wrapper] Before serving request method: %v", req.Method())
		err := fn(ctx, req, rsp)
		log.Log("[Log Wrapper] After serving request")
		return err
	}
}