package logwrapper

import (
	"time"
	"github.com/micro/go-log"
	"context"
	"github.com/micro/go-micro/server"
)

// LogWrapper is a log wrapper
func LogWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func (ctx context.Context, req server.Request, rsp interface{}) error {
		log.Log("[Log Wrapper] Before serving request service: %v", req.Service())
		log.Log("[Log Wrapper] Before serving request method: %v", req.Method())
		log.Log("[Log Wrapper] Before serving request: %#v", req)
		log.Log("[Log Wrapper] Before [%s]", time.Now().Format("02/Jan/2006:15:04:05 -0700"))
		err := fn(ctx, req, rsp)
		log.Log("[Log Wrapper] After [%s]", time.Now().Format("02/Jan/2006:15:04:05 -0700"))
		log.Log("[Log Wrapper] After serving rsp: %#v", rsp)
		return err
	}
}