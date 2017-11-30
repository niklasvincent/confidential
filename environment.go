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
	allowedCharactersRegexp, err := regexp.Compile("[^-./_a-zA-Z0-9]+")
	if err != nil {
		return nil, err
	}
	nonAlphaNumericalCharactersRegexp, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return nil, err
	}

	nameWithOnlyAllowedCharacters := allowedCharactersRegexp.ReplaceAllString(name, "")
	nameWithOutSpecialCharacters := nonAlphaNumericalCharactersRegexp.ReplaceAllString(nameWithOnlyAllowedCharacters, "_")
	nameUpperCased := strings.ToUpper(nameWithOutSpecialCharacters)

	return &nameUpperCased, nil
}
