package imageweb

import (
	"github.com/micro/cli"
	"github.com/micro/go-web"
	"github.com/micro/go-log"
	"github.com/gorilla/mux"

	"time"
)

// Commands add auth api command
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:   "image",
			Usage:  "Run image web",
			Action: aweb,
		},
	}
}

func aweb(ctx *cli.Context) {

	service := web.NewService(
		web.Name("go.micro.web.image"),
		web.RegisterTTL(
			time.Duration(ctx.GlobalInt("register_ttl"))*time.Second,
		),
		web.RegisterInterval(
			time.Duration(ctx.GlobalInt("register_interval"))*time.Second,
		),
	)

	r := mux.NewRouter()
	r.Path("/image")
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/", UploadHandler).Methods("POST")
	r.HandleFunc("/{imgid}", DownloadHandler).Methods("GET")

	service.Handle("/", r)

	LoadConf()

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

}
