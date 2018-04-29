package handler

import (
	"context"

	"btdxcx.com/os/custom-error"

	"github.com/micro/go-log"

	proto "btdxcx.com/micro/taxons-srv/proto/taxons"
)

const (
	svrName = "btdxcx.com/micro/taxons-srv"
)

// Handler is a taxons handler
type Handler struct{}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (h *Handler) Stream(ctx context.Context, req *proto.StreamingRequest, stream proto.Taxons_StreamStream) error {
	log.Logf("Received Taxons.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&proto.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (h *Handler) PingPong(ctx context.Context, stream proto.Taxons_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&proto.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}

func validateCode(code string, method string) error {

	if len(code) < 5 {
		return customerror.BadRequest(svrName, method, "invalid code")
	}

	return nil
}
