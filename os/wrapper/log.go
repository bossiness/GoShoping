package logwrapper

import (
	"fmt"
	"github.com/micro/go-log"
	"context"
	"github.com/micro/go-micro/server"
)

// LogWrapper is a log wrapper
func LogWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func (ctx context.Context, req server.Request, rsp interface{}) error {
		log.Log(fmt.Sprintf("[Log Wrapper] Before serving request service: %v", req.Service()))
		log.Log(fmt.Sprintf("[Log Wrapper] Before serving request method: %v", req.Method()))
		log.Log(fmt.Sprintf("[Log Wrapper] Before serving request: %#v", req.Request()))
		err := fn(ctx, req, rsp)
		log.Log(fmt.Sprintf("[Log Wrapper] After serving rsp: %#v", rsp))
		return err
	}
}