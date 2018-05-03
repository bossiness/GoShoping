package orderapi

import (
	"github.com/micro/cli"
	proto "btdxcx.com/micro/order-srv/proto/order"
)

const (
	clientName     = "com.btdxcx.micro.srv.order"
)

var (
	orderCl proto.OrderClient
)

// API is APIs
type API struct{}

// Commands add command
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:   "cart",
			Usage:  "Run Cart api",
			Action: cart,
		},
		{
			Name:   "checkout",
			Usage:  "Run Checkout api",
			Action: checkout,
		},
		{
			Name:   "order",
			Usage:  "Run Order api",
			Action: order,
		},
	}
}