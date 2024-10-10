package paramstore

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

type ParamStore struct {
	SSMClient *ssm.Client
	SSMPath   string
}

func NewParamStore(ssmPath string) (*ParamStore, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS SDK config, %v", err)
	}

	return &ParamStore{
		SSMClient: ssm.NewFromConfig(cfg),
		SSMPath:   ssmPath,
	}, nil
}

func (p *ParamStore) FormatParamName(name string) string {
	return fmt.Sprintf("%s/%s", p.SSMPath, name)
}

func (p *ParamStore) BuildPutParamInput(name, value string, overwrite bool) *ssm.PutParameterInput {
	return &ssm.PutParameterInput{
		Name:      aws.String(p.FormatParamName(name)),
		Value:     aws.String(value),
		Type:      types.ParameterTypeSecureString,
		Overwrite: aws.Bool(overwrite),
	}
}

func (p *ParamStore) BuildGetParamsByPathInput(next string) *ssm.GetParametersByPathInput {
	return &ssm.GetParametersByPathInput{
		Path:           aws.String(p.SSMPath),
		WithDecryption: aws.Bool(true),
		NextToken:      aws.String(next),
		MaxResults:     aws.Int32(10),
	}
}

func (p *ParamStore) PutParameters(params map[string]string, overwrite bool) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for k, v := range params {
		params := p.BuildPutParamInput(k, v, overwrite)
		r, err := p.SSMClient.PutParameter(ctx, params)
		if err != nil {
			return fmt.Errorf("Error putting parameter %s: %s", k, err)
		}
		fmt.Printf("Parameter added: %s Version: %d\n", *params.Name, r.Version)
	}
	return nil
}

func (p *ParamStore) GetParameters() (map[string]string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	next := ""
	params := make(map[string]string)

	for {
		input := p.BuildGetParamsByPathInput(next)
		result, err := p.SSMClient.GetParametersByPath(ctx, input)

		if err != nil {
			return nil, fmt.Errorf("Error getting parameters: %s", err)
		}

		for _, param := range result.Parameters {
			n := p.ParseParameterName(*param.Name)
			params[n] = *param.Value
		}

		if result.NextToken == nil {
			break
		}
		next = *result.NextToken
	}
	return params, nil
}

func (p *ParamStore) ParseParameterName(name string) string {
	parts := strings.Split(name, "/")
	return parts[len(parts)-1]
}

func FormatParamsAsEnv(params map[string]string) []string {
	var env []string

	for k, v := range params {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	return env
}
