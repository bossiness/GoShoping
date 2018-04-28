package memberapi

import (
	"github.com/micro/cli"
)

// Commands add command
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:   "customer",
			Usage:  "Run customer api",
			Action: cutomersAPIs,
		},
		{
			Name:   "adminuser",
			Usage:  "Run adminuser api",
			Action: adminuserAPIs,
		},
	}
}

const (
	clientName = "com.btdxcx.micro.srv.member"
)

// API is APIs
type API struct{}
