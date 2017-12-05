package environment

import (
	"regexp"
	"strings"
	"fmt"
	"sort"

	"github.com/nlindblad/confidential/aws"
)

// Environment variable
type Variable struct {
	Name  string
	Value string
}

// List of environment variables
type Environment struct {
	variables map[string]string
}

func NewEnvironment() (*Environment) {
	var environment Environment
	environment.variables = make(map[string]string)

	return &environment
}

// Create new list of environment variables from a list of decrypted Amazon SSM parameters
func NewEnvironmentFromDecryptedParameters(decryptedParameters aws.DecryptedParameters) (*Environment, error) {
	var environment Environment
	environment.variables = make(map[string]string)

	for _, param := range decryptedParameters {
		name, err := toEnvironmentVariableName(param.Name)

		if err != nil {
			return nil, err
		}

		environment.Add(Variable{*name, param.Value})
	}

	return &environment, nil
}

// Create a new environment from the runtime environment
func NewEnvironmentFromRuntime(runtimeEnvironment []string) (*Environment, error) {
	var environment Environment
	environment.variables = make(map[string]string)

	for _, variable := range runtimeEnvironment {
		parts := strings.Split(variable, "=")
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid environment variable from runtime: %s", variable)
		}
		name := parts[0]
		value := parts[1]
		environment.Add(Variable{name, value})
	}

	return &environment, nil
}

// Get the union of two sets of environment variable
func (env *Environment) Union(otherEnv *Environment) (*Environment) {
	var newEnvironment Environment
	newEnvironment.variables = make(map[string]string)

	for k, v := range env.variables {
		newEnvironment.variables[k] = v
	}

	for k, v := range otherEnv.variables {
		newEnvironment.variables[k] = v
	}

	return &newEnvironment
}

// Unset a list of environment variables
func (env *Environment) Unset(names []string) {
	for _, name := range names {
		delete(env.variables, name)
	}
}

// Add a new environment variable
func (env *Environment) Add(environmentVariable Variable) {
	env.variables[environmentVariable.Name] = environmentVariable.Value
}

// Get a sorted list representation of the environment
func (env *Environment) AsList() []Variable {
	var environmentList []Variable

	var keys []string
	for k := range env.variables {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, name := range keys {
		value := env.variables[name]
		environmentList = append(environmentList, Variable{name, value})
	}

	return environmentList
}

// Get string array representation of the environment
func (env *Environment) AsStrings() []string {
	var environmentVariables []string

	for _, environmentVariable := range env.AsList() {
		environmentVariables = append(environmentVariables, fmt.Sprintf("%s=%s", environmentVariable.Name, environmentVariable.Value))
	}

	return environmentVariables
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
