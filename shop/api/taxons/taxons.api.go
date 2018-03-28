package taxons

import (
	srvName "btdxcx.com/shop/srv/name"
	"time"
	"github.com/micro/cli"
	"context"
	"log"
	"github.com/micro/go-web"
	"github.com/emicklei/go-restful"
	"github.com/micro/go-micro/client"

	taxons "btdxcx.com/shop/srv/taxons/proto"
)

type TaxonsAPI struct{}

var (
	cl taxons.TaxonsClient
)

const (
	serviceName = "btdxcx.shop.api.taxons"
	clientName = srvName.Taxons
)

func (t *TaxonsAPI) Anything(req *restful.Request, rsp *restful.Response)  {
	log.Print("Received Taxons.Anything API request")
	shopID := req.HeaderParameter("SHOP-ID")
	response, err := cl.Root(context.TODO(), &taxons.TaxonsRequest {
		ShopID: shopID,
	})

	if err != nil {
		rsp.WriteError(500, err)
	}

	rsp.WriteEntity(response)
}

func api(ctx *cli.Context)  {
	service := web.NewService(
		web.Name(serviceName),
		web.RegisterTTL(
			time.Duration(ctx.GlobalInt("register_ttl")) * time.Second,
		),
		web.RegisterInterval(
			time.Duration(ctx.GlobalInt("register_interval")) * time.Second,
		),
	)

//	service.Init()

	cl = taxons.NewTaxonsClient(clientName, client.DefaultClient)

	say := new(TaxonsAPI)
	ws := new(restful.WebService)
	wc := restful.NewContainer()
	ws.Consumes(restful.MIME_XML, restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON, restful.MIME_XML)
	ws.Path("/taxons")
	ws.Route(ws.GET("/").To(say.Anything))
	wc.Add(ws)

	service.Handle("/", wc)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

}


func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "taxons",
			Usage:       "Run taxons api",
			Action: api,
		},
	}
}

// func main()  {
	
// 	service := web.NewService(
// 		web.Name("btdxcx.shop.api.taxons"),
// 	)

// 	service.Init()

// 	cl = taxons.NewTaxonsClient("btdxcx.shop.taxons", client.DefaultClient)

// 	say := new(TaxonsAPI)
// 	ws := new(restful.WebService)
// 	wc := restful.NewContainer()
// 	ws.Consumes(restful.MIME_XML, restful.MIME_JSON)
// 	ws.Produces(restful.MIME_JSON, restful.MIME_XML)
// 	ws.Path("/taxons")
// 	ws.Route(ws.GET("/").To(say.Anything))
// 	wc.Add(ws)

// 	service.Handle("/", wc)

// 	if err := service.Run(); err != nil {
// 		log.Fatal(err)
// 	}
// }