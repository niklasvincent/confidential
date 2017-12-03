package environment

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

func Test_Union_NoOverlap(t *testing.T) {
	expectedEnvironment := []string{
		"A=1",
		"B=2",
		"C=3",
		"D=4",
		"E=5",
		"F=6",
	}

	firstEnvironment := NewEnvironment()
	firstEnvironment.Add(Variable{"A", "1"})
	firstEnvironment.Add(Variable{"B", "2"})
	firstEnvironment.Add(Variable{"D", "4"})
	secondEnvironment := NewEnvironment()
	secondEnvironment.Add(Variable{"C", "3"})
	secondEnvironment.Add(Variable{"E", "5"})
	secondEnvironment.Add(Variable{"F", "6"})

	firstUnion := firstEnvironment.Union(secondEnvironment)

	assert.Equal(t, expectedEnvironment, firstUnion.AsStrings())
}

func Test_Union_HasOverlap(t *testing.T) {
	expectedEnvironment := []string{
		"A=1",
		"B=12",
		"C=13",
		"D=4",
	}

	firstEnvironment := NewEnvironment()
	firstEnvironment.Add(Variable{"A", "1"})
	firstEnvironment.Add(Variable{"B", "2"})
	firstEnvironment.Add(Variable{"C", "3"})
	secondEnvironment := NewEnvironment()
	secondEnvironment.Add(Variable{"B", "12"})
	secondEnvironment.Add(Variable{"C", "13"})
	secondEnvironment.Add(Variable{"D", "4"})

	firstUnion := firstEnvironment.Union(secondEnvironment)

	assert.Equal(t, expectedEnvironment, firstUnion.AsStrings())
}