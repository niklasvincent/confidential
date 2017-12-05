package aws

import (
	"testing"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/stretchr/testify/assert"
	"github.com/aws/aws-sdk-go/aws"
	"strings"
)

type mockSsmParameterStore struct {
	parameters map[string]string
}

func newMockSsmParameterStore() (*mockSsmParameterStore) {
	parameters := make(map[string]string)
	return &mockSsmParameterStore{parameters}
}

func (msp *mockSsmParameterStore) set(key string, value string) {
	msp.parameters[key] = value
}

func (msp *mockSsmParameterStore) get(key string) (*ssm.Parameter) {
	return &ssm.Parameter{Name: aws.String(key), Value: aws.String(msp.parameters[key])}
}

func (msp *mockSsmParameterStore) getParametersByPath(path *string) ([]*ssm.Parameter) {
	var parameters []*ssm.Parameter
	for key := range msp.parameters {
		if strings.HasPrefix(key, *path) {
			parameters = append(parameters, msp.get(key))
		}
	}

	return parameters
}

func TestSsmClient_WithPrefix(t *testing.T) {
	parameterStore := newMockSsmParameterStore()
	parameterStore.set("/test/prod/secret_key", "value1")
	parameterStore.set("/test/prod/database_password", "value2")
	parameterStore.set("/test/prod/cookie_secret", "value3")
	parameterStore.set("/test/dev/cookie_secret", "value4")
	parameterStore.set("/test/dev/cookie_secret", "value5")

	expectedParameters := make(map[string]string)
	expectedParameters["secret_key"] = *parameterStore.get("/test/prod/secret_key").Value
	expectedParameters["database_password"] = *parameterStore.get("/test/prod/database_password").Value
	expectedParameters["cookie_secret"] = *parameterStore.get("/test/prod/cookie_secret").Value

	client := SsmClient{ssm.New(session.Must(session.NewSession()))}

	client.client.Handlers.Clear()
	client.client.Handlers.Send.PushBack(func(r *request.Request) {
		getParameterByPathInput, inputOk := r.Params.(*ssm.GetParametersByPathInput)
		getParameterByPathOutput, outputOk := r.Data.(*ssm.GetParametersByPathOutput)
		if inputOk && outputOk {
			path := getParameterByPathInput.Path
			getParameterByPathOutput.Parameters = parameterStore.getParametersByPath(path)
		}
	})

	retrievedParameters, _ := client.WithPrefix("/test/prod")

	assert.Equal(t, len(expectedParameters), len(retrievedParameters))
	for _, retrievedParameter := range retrievedParameters {
		assert.Equal(t, expectedParameters[retrievedParameter.Name], retrievedParameter.Value)
	}
}