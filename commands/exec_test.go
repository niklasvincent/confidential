package commands

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func Test_extractCommandWithFlags(t *testing.T) {
	var err error

	cmd, err := extractCommandWithFlags([]string{"--prefix", "/test/prod", "--region", "eu-west-1", "exec", "--", "env"})
	assert.Equal(t, "env", cmd.executable)
	assert.Equal(t, 0, len(cmd.arguments))

	cmd, err = extractCommandWithFlags([]string{"--prefix", "/test/prod", "--region", "eu-west-1", "exec", "--", "echo", "\"${SECRET_KEY}\""})
	assert.Equal(t, "echo", cmd.executable)
	assert.Equal(t, 1, len(cmd.arguments))

	cmd, err = extractCommandWithFlags([]string{"--prefix", "/test/prod", "--region", "eu-west-1", "exec", "--"})
	assert.NotNil(t, err)
}