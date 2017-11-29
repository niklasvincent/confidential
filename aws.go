package main

import (
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

type ssmClient struct {
	client *ssm.SSM
}

type parameter struct {
	Name     string
	Value	 string
}

type parameters []parameter

func NewClient() *ssmClient {
	session := session.Must(session.NewSession())
	return &ssmClient{ssm.New(session) }
}

func (s ssmClient) ParamList(filter string) (*ssm.DescribeParametersOutput, error) {
	// TODO(nlindblad): Add support for pagination
	params := &ssm.DescribeParametersInput{
		MaxResults: aws.Int64(50),
		Filters: []*ssm.ParametersFilter{
			{
				Values: []*string{
					aws.String(filter),
				},
				Key: aws.String("Name"),
			},
		},
	}

	return s.client.DescribeParameters(params)
}

func (s ssmClient) WithPrefix(prefix string) (parameters, error) {
	var parameters parameters
	resp, err := s.ParamList(prefix)
	if err != nil {
		return nil, err
	}
	for _, param := range resp.Parameters {
		parametersInput := &ssm.GetParametersInput{
			Names: []*string{param.Name},
			WithDecryption: aws.Bool(true),
		}

		r, err := s.client.GetParameters(parametersInput)
		if err != nil {
			return nil, err
		}

		nameWithoutPrefix := string([]rune(*param.Name)[len(prefix) + 1:])

		parameters = append(parameters, parameter{nameWithoutPrefix, *r.Parameters[0].Value})
	}
	return parameters, nil
}