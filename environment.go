package main

import (
	"regexp"
	"strings"
)

// Environment variable
type EnvironmentVariable struct {
	Name  string
	Value string
}

// List of environment variables
type Environment []EnvironmentVariable

// Create new list of environment variables from a list of decrypted Amazon SSM parameters
func NewEnvironment(decryptedParameters DecryptedParameters) (Environment, error) {
	var environment Environment

	for _, param := range decryptedParameters {
		name, err := toEnvironmentVariableName(param.Name)

		if err != nil {
			return nil, err
		}

		environment = append(environment, EnvironmentVariable{*name, param.Value})
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
