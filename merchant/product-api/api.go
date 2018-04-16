package productapi

import (
	"github.com/micro/cli"
)

// API is APIs
type API struct{}

// Commands add command
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:   "product",
			Usage:  "Run product api",
			Action: api,
		},
		{
			Name:   "products",
			Usage:  "Run products api",
			Action: apis,
		},
	}
}
