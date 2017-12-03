package commands

import (
	"fmt"

	"github.com/urfave/cli"

	"github.com/nlindblad/confidential/environment"
	"github.com/nlindblad/confidential/aws"
)

// Get a mandatory command line flag and error if it does not exist
func GetMandatoryFlag(c *cli.Context, name string) (*string, error) {
	value := ""
	if c.GlobalIsSet(name) {
		value = c.GlobalString(name)
	} else {
		value = c.String(name)
	}

	if value == "" {
		return nil, fmt.Errorf("required flag --%s not provided", name)
	}

	return &value, nil
}

func EnvironmentForAwsCredentialsProfile(profile string) (*environment.Environment, error) {
	environmentWithAwsCredentials := environment.NewEnvironment()

	credentials, err:= aws.GetAwsCredentialsForProfile(profile)
	if err != nil {
		return nil, err
	}

	value, err := credentials.Get()
	if err != nil {
		return nil, err
	}

	environmentWithAwsCredentials.Add(environment.Variable{"AWS_ACCESS_KEY_ID", value.AccessKeyID})
	environmentWithAwsCredentials.Add(environment.Variable{"AWS_SECRET_ACCESS_KEY", value.SecretAccessKey})
	environmentWithAwsCredentials.Add(environment.Variable{"AWS_SESSION_TOKEN", value.SessionToken})

	return environmentWithAwsCredentials, nil
}

// Create new wrapped Amazon SSM client from CLI context
func NewClientFromContext(c *cli.Context) (*aws.SsmClient, error) {
	region, err := GetMandatoryFlag(c, "region")
	if err != nil {
		return nil, err
	}

	if c.GlobalIsSet("profile") {
		profile := c.GlobalString("profile")
		credentials, err := aws.GetAwsCredentialsForProfile(profile)
		if err != nil {
			return nil, err
		}
		return aws.NewClientWithCredentials(*region, credentials)
	} else {
		return aws.NewClient(*region)
	}
}

// Retrieve Amazon SSM parameters as environment variables for a given prefix from CLI context
func RetrieveEnvironmentVariablesFromContext(c *cli.Context) (*environment.Environment, error) {
	prefix, err := GetMandatoryFlag(c, "prefix")
	client, err := NewClientFromContext(c)
	if err != nil {
		return nil, err
	}

	parameters, err := client.WithPrefix(*prefix)
	if err != nil {
		return nil, err
	}

	environment, err := environment.NewEnvironmentFromDecryptedParameters(parameters)
	if err != nil {
		return nil, err
	}

	if c.GlobalIsSet("forwarded-profile") {
		environmentWithAwsCredentials, err := EnvironmentForAwsCredentialsProfile(c.GlobalString("forwarded-profile"))
		if err != nil {
			return nil, err
		}
		environment = environment.Union(environmentWithAwsCredentials)
	}

	return environment, nil
}


func GetCommands() []cli.Command {
	return []cli.Command{
		{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "retrieve and atomically output environment variables to a file",
			Action:  output,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "env-file",
					Usage: "Output file to write environment variables to",
				},
			},
		},
		{
			Name:    "exec",
			Aliases: []string{"e"},
			Usage:   "retrieve environment variables and execute command with an updated environment",
			Action:  execRun,
		},
	}
}