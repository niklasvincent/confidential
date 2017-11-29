package main

import (
	"regexp"
	"strings"
)

type environmentVariable struct {
	Name  string
	Value string
}

type environment []environmentVariable

func NewEnvironment(parameters parameters) (environment, error) {
	var environment environment

	for _, param := range parameters {
		name, err := toEnvironmentVariableName(param.Name)

		if err != nil {
			return nil, err
		}

		environment = append(environment, environmentVariable{*name, param.Value})
	}

	return environment, nil
}

func toEnvironmentVariableName(name string) (*string, error) {
	reg, err := regexp.Compile("[^_a-zA-Z0-9]+")
	if err != nil {
		return nil, err
	}

	nameWithoutSlashes := strings.Replace(name, "/", "_", -1)
	nameWithOnlyAlphanumericalCharacters := reg.ReplaceAllString(nameWithoutSlashes, "")
	nameUpperCased := strings.ToUpper(nameWithOnlyAlphanumericalCharacters)

	return &nameUpperCased, nil
}
