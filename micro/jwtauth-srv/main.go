package main

import (
	"btdxcx.com/micro/jwtauth-srv/db"
	"btdxcx.com/micro/jwtauth-srv/db/mongodb"
	"github.com/micro/cli"
	"time"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"btdxcx.com/micro/jwtauth-srv/handler"

	proto "btdxcx.com/micro/jwtauth-srv/proto/auth"

)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("com.btdxcx.micro.srv.jwtauth"),
		micro.Version("latest"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*15),
		micro.Metadata(map[string]string{
			"type": "jwtauth srv",
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
	proto.RegisterJwtAuthHandler(service.Server(), new(handler.Handler))

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
