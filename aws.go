package main

import (
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

type ssmClient struct {
	client *ssm.SSM
}

type DecryptedParameter struct {
	Name  string
	Value string
}

type DecryptedParameters []DecryptedParameter

func NewClient() *ssmClient {
	session := session.Must(session.NewSession())
	return &ssmClient{ssm.New(session)}
}

func (s ssmClient) paramListPaginated(prefix string, nextToken *string) ([]ssm.Parameter, *string, error) {
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

func (s ssmClient) ParamList(prefix string) (*[]ssm.Parameter, error) {
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

func (s ssmClient) WithPrefix(prefix string) (DecryptedParameters, error) {
	var parameters DecryptedParameters

	retrievedParameters, err := s.ParamList(prefix)
	if err != nil {
		return nil, err
	}

	for _, parameter := range *retrievedParameters {
		nameWithoutPrefix := string([]rune(*parameter.Name)[len(prefix)+1:])
		parameters = append(parameters, DecryptedParameter{nameWithoutPrefix, *parameter.Value})
	}

	return parameters, nil
}
