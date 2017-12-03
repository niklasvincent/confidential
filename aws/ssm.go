package aws

import (
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

// Wrapped Amazon SSM client
type SsmClient struct {
	client *ssm.SSM
}

// Decrypted parameter (without prefix) from Amazon SSM
type DecryptedParameter struct {
	Name  string
	Value string
}

// List of decrypted parameters (without prefix) from Amazon SSM
type DecryptedParameters []DecryptedParameter

// Create new wrapped Amazon SSM client
func NewClient(awsRegion string) (*SsmClient, error) {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil {
		return nil, err
	}

	return &SsmClient{ssm.New(session)}, nil
}

// Create new wrapped Amazon SSM client
func NewClientWithCredentials(awsRegion string, credentials *credentials.Credentials) (*SsmClient, error) {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
		Credentials: credentials,
	})
	if err != nil {
		return nil, err
	}

	return &SsmClient{ssm.New(session)}, nil
}

func (s SsmClient) paramListPaginated(prefix string, nextToken *string) ([]ssm.Parameter, *string, error) {
	var parameters []ssm.Parameter

	getParametersByPathInput := &ssm.GetParametersByPathInput{
		NextToken:  nextToken,
		Path: &prefix,
		Recursive: aws.Bool(true),
		WithDecryption: aws.Bool(true),
	}

	result, err := s.client.GetParametersByPath(getParametersByPathInput)
	if err != nil {
		return nil, nil, err
	}

	for _, parameter := range result.Parameters {
		parameters = append(parameters, *parameter)
	}

	return parameters, result.NextToken, nil
}

func (s SsmClient) paramList(prefix string) (*[]ssm.Parameter, error) {
	parameters, nextToken, err := s.paramListPaginated(prefix, nil)
	if err != nil {
		return nil, err
	}

	for nextToken != nil {
		var additionalParameters []ssm.Parameter
		additionalParameters, nextToken, err = s.paramListPaginated(prefix, nextToken)
		if err != nil {
			return nil, err
		}
		for _, parameter := range additionalParameters {
			parameters = append(parameters, parameter)
		}
	}

	return &parameters, nil
}

func (s SsmClient) WithPrefix(prefix string) (DecryptedParameters, error) {
	var parameters DecryptedParameters

	retrievedParameters, err := s.paramList(prefix)
	if err != nil {
		return nil, err
	}

	for _, parameter := range *retrievedParameters {
		nameWithoutPrefix := string([]rune(*parameter.Name)[len(prefix)+1:])
		parameters = append(parameters, DecryptedParameter{nameWithoutPrefix, *parameter.Value})
	}

	return parameters, nil
}