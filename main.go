package main

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/supermicah/dionysus-admin/cmd"
)

// VERSION Usage: go build -ldflags "-X main.VERSION=x.x.x"
var VERSION = "v1.0.0"

// @title dionysus-admin
// @version v1.0.0
// @description A test API service based on golang.
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @schemes http https
// @basePath /
func main() {
	app := cli.NewApp()
	app.Name = "dionysus-admin"
	app.Version = VERSION
	app.Usage = "A test API service based on golang."
	app.Commands = []*cli.Command{
		cmd.StartCmd(),
		cmd.StopCmd(),
		cmd.VersionCmd(VERSION),
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
