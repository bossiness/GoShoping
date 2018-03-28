package main

import (
	"btdxcx.com/shop/srv/taxons"
	"os"

	ccli "github.com/micro/cli"
	"github.com/micro/go-micro/cmd"
)

func init() {
	os.Setenv("SHOP_SVR", "svr")
}

func setup(app *ccli.App) {
	// common flags
	app.Flags = append(app.Flags) // ccli.StringFlag{
	// 	Name: "register_ttl", 
	// 	EnvVar: "MICRO_REGISTER_TTL",
	// 	Usage:  "Register TTL in seconds",
	// },
	// ccli.IntFlag{
	// 	Name:   "register_interval",
	// 	EnvVar: "MICRO_REGISTER_INTERVAL",
	// 	Usage:  "Register interval in seconds",
	// },

}

func main() {
	app := cmd.App()

	app.Commands = append(app.Commands, srvtaxons.Commands()...)
	app.Action = func(context *ccli.Context) {
		ccli.ShowAppHelp(context)
	}

	setup(app)

	cmd.Init(
		cmd.Name("svr"),
		cmd.Description("Shop SVRs"),
		cmd.Version("latest"),
	)
}
