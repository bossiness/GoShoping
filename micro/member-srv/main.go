package main

import (
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"btdxcx.com/micro/member-srv/handler"

	proto "btdxcx.com/micro/member-srv/proto/member"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("com.btdxcx.micro.srv.member"),
		micro.Version("v1"),
	)

	// Register Handler
	proto.RegisterCustomerHandler(service.Server(), new(handler.Handler))

	// Initialise service
	service.Init()

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
