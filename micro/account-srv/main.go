package main

import (
	"github.com/micro/go-micro/server"
	"btdxcx.com/os/wrapper"
	"time"

	"btdxcx.com/micro/account-srv/db"

	"btdxcx.com/micro/account-srv/db/mongodb"
	"btdxcx.com/micro/account-srv/handler"
	"github.com/micro/cli"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"

	proto "btdxcx.com/micro/account-srv/proto/account"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Server(
			server.NewServer(server.WrapHandler(logwrapper.LogWrapper)),
		),
		micro.Name("com.btdxcx.micro.srv.account"),
		micro.Version("v1"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*15),
		micro.Metadata(map[string]string{
			"type": "taxons srv",
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
	proto.RegisterAccountHandler(service.Server(), new(handler.Handler))

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
