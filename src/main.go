package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gocker/src/command"
	"os"
)

const usage = `Gocker is a simple container runtime implementation.`

func main() {
	app := cli.NewApp()
	app.Name = "Gocker"
	app.Usage = usage
	app.Commands = []cli.Command{
		command.InitGockerCMD,
		command.RunGockerCMD,
	}
	app.Before = func(context *cli.Context) error {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetOutput(os.Stdout)
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
