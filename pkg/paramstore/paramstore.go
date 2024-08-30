package paramstore

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/joho/godotenv"
)

// FileExists checks if a file exists and is not a directory.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func LoadEnvFile(envPath string) (map[string]string, error) {

	if !FileExists(envPath) {
		return nil, fmt.Errorf("Environment file not found at %s \n", envPath)
	}

	env, err := godotenv.Read(envPath)
	if err != nil {
		return nil, fmt.Errorf("Error reading environment file: %s \n", err)
	}

	return env, nil
}

func getSSMClient() (*ssm.Client, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	return ssm.NewFromConfig(cfg), nil
}

func BuildSSMPutParameterInput(name, value string, overwrite bool) ssm.PutParameterInput {
	return ssm.PutParameterInput{
		Name:      aws.String(name),
		Value:     aws.String(value),
		Type:      types.ParameterTypeSecureString,
		Overwrite: aws.Bool(overwrite),
	}
}

func BuildParameterStoreVariablePath(path, name, value string) string {
	return fmt.Sprintf("%s/%s", path, name)
}

func PutParametersFromEnvFile(env map[string]string, path string, overwrite bool) error {

	client, err := getSSMClient()
	if err != nil {
		return err
	}

	for k, v := range env {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		varPath := BuildParameterStoreVariablePath(path, k, v)
		params := BuildSSMPutParameterInput(varPath, v, overwrite)
		result, err := client.PutParameter(ctx, &params)

		if err != nil {
			return fmt.Errorf("Error putting parameter %s: %s", k, err)
		}
		fmt.Printf("Parameter added: %s Version: %d\n", varPath, result.Version)
	}
	return nil
}
