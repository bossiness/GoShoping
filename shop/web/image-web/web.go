package imageweb

import (
	"github.com/micro/cli"
	"github.com/micro/go-web"
	"github.com/micro/go-log"

	"time"
	"net/http"
)

// Commands add auth api command
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:   "image",
			Usage:  "Run image web",
			Action: action,
		},
	}
}

func action(ctx *cli.Context) {

	service := web.NewService(
		web.Name("go.micro.web.image"),
		web.RegisterTTL(
			time.Duration(ctx.GlobalInt("register_ttl"))*time.Second,
		),
		web.RegisterInterval(
			time.Duration(ctx.GlobalInt("register_interval"))*time.Second,
		),
	)

	service.HandleFunc("/{imgid}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			DownloadHandler(w, r)
			return
		}
	})
	service.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			UploadHandler(w, r)
			return
		}
		HomeHandler(w, r)
	})

	LoadConf()

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

}
