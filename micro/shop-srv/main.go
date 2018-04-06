package main

import (
	"time"

	"btdxcx.com/micro/shop-srv/db"
	"btdxcx.com/micro/shop-srv/db/mongodb"
	"btdxcx.com/micro/shop-srv/handler"
	"github.com/micro/cli"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"

	proto "btdxcx.com/micro/shop-srv/proto/shop"
	dproto "btdxcx.com/micro/shop-srv/proto/shop/details"
)

func main() {

	// New Service
	service := micro.NewService(
		micro.Name("com.btdxcx.micro.srv.shop"),
		micro.Version("latest"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*15),
		micro.Metadata(map[string]string{
			"type": "shop srv",
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
	keyHandler := &handler.KeyHandler{ Tags: []string{ "back", "mini" } }
	proto.RegisterShopKeyHandler(service.Server(), keyHandler)

	detailsHandler := &handler.DetailsHandler{}
	dproto.RegisterShopHandler(service.Server(), detailsHandler)

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
