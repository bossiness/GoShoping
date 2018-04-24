package main

import (
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"btdxcx.com/micro/order-srv/handler"

	proto "btdxcx.com/micro/order-srv/proto/order"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("com.btdxcx.micro.srv.order"),
		micro.Version("latest"),
	)

	// Register Handler
	proto.RegisterExampleHandler(service.Server(), new(handler.Example))

	// Initialise service
	service.Init()

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
