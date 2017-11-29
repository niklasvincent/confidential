package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
	"github.com/dchest/safefile"
)

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "prefix",
			Usage: "parameter prefix",
		},
		cli.StringFlag{
			Name: "env-file",
			Usage: "output environment file",
		},
	}

	app.Name = "confidential"
	app.Authors = []cli.Author{
		{
			Name:  "Niklas Lindblad",
			Email: "niklas@lindblad.info",
		},
	}
	app.Action = func(c *cli.Context) error {
		prefix := c.String("prefix")
		if prefix == "" {
			return fmt.Errorf("no prefix specififed")
		}

		envFile := c.String("env-file")
		if envFile == "" {
			return fmt.Errorf("no env-file specififed")
		}

		client := NewClient()

		parameters, err := client.WithPrefix(prefix)
		if err != nil {
			return err
		}

		environment, err := NewEnvironment(parameters)
		if err != nil {
			return err
		}

		f, err := safefile.Create(envFile, 0644)
		if err != nil {
			return nil
		}
		defer f.Close()

		for _, environmentVariable := range environment {
			_, err := f.WriteString(fmt.Sprintf("%s=%s\n", environmentVariable.Name, environmentVariable.Value))
			if err != nil {
				return err
			}
		}

		err = f.Commit()
		if err != nil {
			return err
		}

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}