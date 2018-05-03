package main

import (
	"os"

	"btdxcx.com/applet/product-api"
	"btdxcx.com/applet/shop-api"
	"btdxcx.com/applet/taxons-api"
	"btdxcx.com/applet/weixin-api"
	"btdxcx.com/applet/order-api"

	ccli "github.com/micro/cli"
	"github.com/micro/go-micro/cmd"
)

func init() {
	os.Setenv("APPLET_API", "api")
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
	app.Commands = append(app.Commands, productapi.Commands()...)
	app.Commands = append(app.Commands, taxonsapi.Commands()...)
	app.Commands = append(app.Commands, shopapi.Commands()...)
	app.Commands = append(app.Commands, weixinapi.Commands()...)
	app.Commands = append(app.Commands, orderapi.Commands()...)
	app.Action = func(context *ccli.Context) {
		ccli.ShowAppHelp(context)
	}

	setup(app)

	cmd.Init(
		cmd.Name("applet"),
		cmd.Description("applet apis"),
		cmd.Version("v1"),
	)
}
