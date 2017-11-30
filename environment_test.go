package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func Test_toEnvironmentVariableName(t *testing.T) {
	environmentVariableName, _ := toEnvironmentVariableName("secret_key")
	assert.Equal(t, "SECRET_KEY", *environmentVariableName)

	environmentVariableName, _ = toEnvironmentVariableName("secret/key")
	assert.Equal(t, "SECRET_KEY", *environmentVariableName)

	environmentVariableName, _ = toEnvironmentVariableName("secret key")
	assert.Equal(t, "SECRETKEY", *environmentVariableName)

	environmentVariableName, _ = toEnvironmentVariableName(".secret-key")
	assert.Equal(t, "_SECRET_KEY", *environmentVariableName)

	environmentVariableName, _ = toEnvironmentVariableName(".secret#-@key")
	assert.Equal(t, "_SECRET_KEY", *environmentVariableName)
}
