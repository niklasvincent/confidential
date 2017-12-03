package commands

import (
	"github.com/dchest/safefile"
	"github.com/nlindblad/confidential/environment"
	"github.com/urfave/cli"
	"fmt"
)

func writeEnvironmentToFile(environment *environment.Environment, envFile *string) (error) {
	f, err := safefile.Create(*envFile, 0644)
	if err != nil {
		return nil
	}
	defer f.Close()

	for _, environmentVariable := range environment.AsStrings() {
		_, err := f.WriteString(fmt.Sprintf("%s\n", environmentVariable))
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

func output(c *cli.Context) error {
	environment, err := RetrieveEnvironmentVariablesFromContext(c)
	if err != nil {
		return err
	}

	envFile, err := GetMandatoryFlag(c, "env-file")
	if err != nil {
		return err
	}

	writeEnvironmentToFile(environment, envFile)

	return nil
}