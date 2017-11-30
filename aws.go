package main

import (
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

type ssmClient struct {
	client    *ssm.SSM
}

type parameter struct {
	Name  string
	Value string
}

type parameters []parameter

func NewClient() *ssmClient {
	session := session.Must(session.NewSession())
	return &ssmClient{ssm.New(session)}
}

func (s ssmClient) paramListPaginated(filter string, nextToken *string) ([]ssm.ParameterMetadata, *string, error) {
	var parametersMetadata []ssm.ParameterMetadata

	describeParameterInput := &ssm.DescribeParametersInput{
		MaxResults: aws.Int64(50),
		NextToken:  nextToken,
		Filters: []*ssm.ParametersFilter{
			{
				Values: []*string{
					aws.String(filter),
				},
				Key: aws.String("Name"),
			},
		},
	}

	result, err := s.client.DescribeParameters(describeParameterInput)
	if err != nil {
		return nil, nil, err
	}

	for _, paramMetaData := range result.Parameters {
		parametersMetadata = append(parametersMetadata, *paramMetaData)
	}

	return parametersMetadata, result.NextToken, nil
}

func (s ssmClient) ParamList(filter string) (*[]ssm.ParameterMetadata, error) {
	parametersMetadata, nextToken, err := s.paramListPaginated(filter, nil)
	if err != nil {
		return nil, err
	}

	for nextToken != nil {
		var additionalParametersMetadata []ssm.ParameterMetadata
		additionalParametersMetadata, nextToken, err = s.paramListPaginated(filter, nextToken)
		if err != nil {
			return nil, err
		}
		for _, paramMetaData := range additionalParametersMetadata {
			parametersMetadata = append(parametersMetadata, paramMetaData)
		}
	}

	return &parametersMetadata, nil
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func (s ssmClient) WithPrefix(prefix string) (parameters, error) {
	var parameters parameters

	parameterMetadata, err := s.ParamList(prefix)
	if err != nil {
		return nil, err
	}

	// The Parameter Store API supports requesting 10 parameters at a time.
	// This divides the list of parameter metadata into chunks of 10 in order
	// to make as few API calls as possible.
	for start := 0; start < len(*parameterMetadata); start += 10 {
		end := min(start + 10, len(*parameterMetadata))
		batch := (*parameterMetadata)[start:end]
		parameterNames := make([]*string, len(batch))

		for i, parameterMetadata := range batch {
			parameterNames[i] = parameterMetadata.Name
		}
		parametersInput := &ssm.GetParametersInput{
			Names:          parameterNames,
			WithDecryption: aws.Bool(true),
		}

		result, err := s.client.GetParameters(parametersInput)
		if err != nil {
			return nil, err
		}

		for _, param := range result.Parameters {
			nameWithoutPrefix := string([]rune(*param.Name)[len(prefix)+1:])
			parameters = append(parameters, parameter{nameWithoutPrefix, *param.Value})
		}
	}

	return parameters, nil
}
