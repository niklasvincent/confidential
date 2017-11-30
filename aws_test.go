package main

import (
	"testing"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/stretchr/testify/assert"
)


// Strings are treated as constants, so we cannot get the address of a string literal
// without this hack: https://groups.google.com/forum/#!topic/golang-nuts/mKJbGRRJm7c
func ptrTo(s string) *string {
	return &s;
}

func TestSsmClient_WithPrefix(t *testing.T) {
	testParameters := make(map[string]string)
	testParameters["/test/prod/secret_key"] = "value1"
	testParameters["/test/prod/database_password"] = "value2"
	testParameters["/test/prod/cookie_secret"] = "value3"

	expectedParameters := make(map[string]string)
	expectedParameters["secret_key"] = testParameters["/test/prod/secret_key"]
	expectedParameters["database_password"] = testParameters["/test/prod/database_password"]
	expectedParameters["cookie_secret"] = testParameters["/test/prod/cookie_secret"]

	client := ssmClient{ssm.New(session.Must(session.NewSession()))}

	client.client.Handlers.Clear()
	client.client.Handlers.Send.PushBack(func(r *request.Request) {
		describeParametersOutput, ok := r.Data.(*ssm.DescribeParametersOutput)
		if ok {
			for name, _ := range testParameters {
				describeParametersOutput.Parameters = append(describeParametersOutput.Parameters, &ssm.ParameterMetadata{Name: ptrTo(name), Type: ptrTo("SecretString")})
			}
		}
		getParametersOutput, ok := r.Data.(*ssm.GetParametersOutput)
		if ok {
			for _, name := range r.Params.(*ssm.GetParametersInput).Names {
				getParametersOutput.Parameters = append(getParametersOutput.Parameters, &ssm.Parameter{Name: name, Value: ptrTo(testParameters[*name])})
			}
		}
	})

	retrievedParameters, _ := client.WithPrefix("/test/prod")

	assert.Equal(t, len(expectedParameters), len(retrievedParameters))
	for _, retrievedParameter := range retrievedParameters {
		assert.Equal(t, expectedParameters[retrievedParameter.Name], retrievedParameter.Value)
	}
}