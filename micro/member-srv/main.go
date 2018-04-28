package main

import (
	"btdxcx.com/micro/member-srv/db"
	"time"

	"btdxcx.com/micro/member-srv/db/mongodb"
	"btdxcx.com/micro/member-srv/handler"
	"github.com/micro/cli"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
	"btdxcx.com/os/wrapper"

	proto "btdxcx.com/micro/member-srv/proto/member"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Server(
			server.NewServer(server.WrapHandler(logwrapper.LogWrapper)),
		),
		micro.Name("com.btdxcx.micro.srv.member"),
		micro.Version("v1"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*15),
		micro.Metadata(map[string]string{
			"type": "member srv",
		}),

		micro.Flags(
			cli.StringFlag{
				Name:   "database_url",
				EnvVar: "DATABASE_URL",
				Usage:  "The database URL e.g localhost:27017",
			},
		),

		micro.Action(func(c *cli.Context) {
			if len(c.String("database_url")) > 0 {
				mongodb.DBUrl = c.String("database_url")
			}
		}),
	)

	// Register Handler
	proto.RegisterCustomerHandler(service.Server(), new(handler.CustomerHandler))
	proto.RegisterAdminUserHandler(service.Server(), new(handler.AdminUserHandler))

	// Initialise service
	service.Init()

	// initialise database
	if err := db.Init(); err != nil {
		log.Fatal(err)
	}
	defer db.Deinit()

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
