package main

import (
	// "btdxcx.com/micro/shop-srv/wrapper/inspection/shop-key"
	"time"

	"btdxcx.com/micro/taxons-srv/db"
	"btdxcx.com/micro/taxons-srv/handler"
	"btdxcx.com/micro/taxons-srv/subscriber"
	"btdxcx.com/micro/taxons-srv/wrapper"
	"github.com/micro/cli"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/server"

	"btdxcx.com/micro/taxons-srv/db/mongodb"
	proto "btdxcx.com/micro/taxons-srv/proto/imp"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Server(
			server.NewServer(server.WrapHandler(wrapper.LogWrapper)),
		),
		micro.Name("com.btdxcx.shop.srv.taxons"),
		micro.Version("latest"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*15),
		micro.Metadata(map[string]string{
			"type": "taxons srv",
		}),

		micro.Flags(
			cli.StringFlag{
				Name:   "database_url",
				EnvVar: "DATABASE_URL",
				Usage:  "The database URL e.g root@tcp(127.0.0.1:3306)/auth",
			},
		),

		micro.Action(func(c *cli.Context) {
			if len(c.String("database_url")) > 0 {
				mongodb.Url = c.String("database_url")
			}
		}),

	)


	// Register Handler
	proto.RegisterTaxonsHandler(service.Server(), new(handler.Handler))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("topic.com.btdxcx.shop.srv.taxons", service.Server(), new(subscriber.Receiver))

	// Register Function as Subscriber
	micro.RegisterSubscriber("topic.com.btdxcx.shop.srv.taxons", service.Server(), subscriber.Handler)

	// Initialise service
	service.Init()

	// initialise database
	if err := db.Init(); err != nil {
		log.Fatal(err)
	}

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
