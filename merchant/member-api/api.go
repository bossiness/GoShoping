package memberapi

import (
	"github.com/micro/cli"
)

// API is APIs
type API struct{}

// Commands add command
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:   "customer",
			Usage:  "Run customer api",
			Action: cutomersAPIs,
		},
	}
}
