package srvtaxons

import (
	srvName "btdxcx.com/shop/srv/name"
	"log"
	"github.com/micro/go-micro"
	proto "btdxcx.com/shop/srv/taxons/proto"
	"golang.org/x/net/context"
	"github.com/micro/cli"
)

type TaxonsHandler struct {}

func (h *TaxonsHandler) Root (ctx context.Context, req *proto.TaxonsRequest, rsp *proto.TaxonsResponse) error {
	rsp.Code = "Shop ID:" + req.GetShopID()
	return nil
}

func srv(ctx *cli.Context)  {

	service := micro.NewService(
		micro.Name(srvName.Taxons),
		micro.Version("latest"),
	)

	// service.Init()

	proto.RegisterTaxonsHandler(service.Server(), new(TaxonsHandler))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

}

func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:   "taxons",
			Usage:  "Run taxons svr",
			Action: srv,
		},
	}
}


// func main()  {

// 	service := micro.NewService(
// 		micro.Name(srvName.Taxons),
// 		micro.Version("latest"),
// 	)

// 	service.Init()

// 	proto.RegisterTaxonsHandler(service.Server(), new(TaxonsHandler))

// 	if err := service.Run(); err != nil {
// 		log.Fatal(err)
// 	}
	
// }