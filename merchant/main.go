package main

import (
	"os"

	"btdxcx.com/merchant/member-api"
	"btdxcx.com/merchant/product-api"
	"btdxcx.com/merchant/taxons-api"
	"btdxcx.com/merchant/shop-api"

	ccli "github.com/micro/cli"
	"github.com/micro/go-micro/cmd"
)

func init() {
	os.Setenv("MERCHANT_API", "api")
}

func setup(app *ccli.App) {
	// common flags
	app.Flags = append(app.Flags,
		ccli.StringFlag{
			Name:   "register_ttl",
			EnvVar: "MICRO_REGISTER_TTL",
			Usage:  "Register TTL in seconds",
		},
		ccli.IntFlag{
			Name:   "register_interval",
			EnvVar: "MICRO_REGISTER_INTERVAL",
			Usage:  "Register interval in seconds",
		},
	)
}

func main() {
	app := cmd.App()
	app.Commands = append(app.Commands, taxonsapi.Commands()...)
	app.Commands = append(app.Commands, productapi.Commands()...)
	app.Commands = append(app.Commands, memberapi.Commands()...)
	app.Commands = append(app.Commands, shopapi.Commands()...)
	app.Action = func(context *ccli.Context) {
		ccli.ShowAppHelp(context)
	}

	setup(app)

	cmd.Init(
		cmd.Name("merchant"),
		cmd.Description("merchant apis"),
		cmd.Version("v1"),
	)
}
