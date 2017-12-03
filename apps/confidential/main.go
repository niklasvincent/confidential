package main

import (
	"log"
	"os"
	"sort"

	"github.com/urfave/cli"


	"github.com/nlindblad/confidential/commands"
)

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "region",
			Usage: "AWS region e.g. eu-west-1",
			EnvVar: "AWS_REGION",
		},
		cli.StringFlag{
			Name: "prefix",
			Usage: "Amazon SSM parameter prefix",
		},
		cli.StringFlag{
			Name: "profile",
			Usage: "AWS profile to use when calling Amazon SSM",
			EnvVar: "AWS_PROFILE",
		},
		cli.StringFlag{
			Name: "forwarded-profile",
			Usage: "AWS profile to forward credentials for in the created environment",
			EnvVar: "AWS_FORWARDED_PROFILE",
		},
	}

	app.Name = "confidential"
	app.Authors = []cli.Author{
		{
			Name:  "Niklas Lindblad",
			Email: "niklas@lindblad.info",
		},
	}

	app.Commands = commands.GetCommands()

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}