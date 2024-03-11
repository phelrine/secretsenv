package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type SecretResult map[string]string

type SecretLoader interface {
	Load(secretId string, option SecretOption) (SecretResult, error)
}

type AWSSecretsManagerLoader struct {
}

func NewAWSLoader() *AWSSecretsManagerLoader {
	return &AWSSecretsManagerLoader{}
}

func (l *AWSSecretsManagerLoader) Load(secretId string, option SecretOption) (SecretResult, error) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithAssumeRoleCredentialOptions(func(options *stscreds.AssumeRoleOptions) {
			options.TokenProvider = func() (string, error) {
				var v string
				fmt.Fprintf(os.Stderr, "Assume Role MFA token code: ")
				_, err := fmt.Scanln(&v)
				return v, err
			}
		}),
	)
	if err != nil {
		return nil, err
	}
	client := secretsmanager.NewFromConfig(cfg)
	result, err := client.GetSecretValue(
		ctx,
		&secretsmanager.GetSecretValueInput{
			SecretId: aws.String(secretId),
		},
	)
	if err != nil {
		return nil, err
	}
	var secretMap map[string]interface{}
	err = json.Unmarshal([]byte(*result.SecretString), &secretMap)
	if err != nil {
		return nil, err
	}
	secretResult := make(SecretResult)
	for secretKey, secretValue := range secretMap {
		switch v := secretValue.(type) {
		case string:
			secretResult[secretKey] = v
		case float64:
			secretResult[secretKey] = fmt.Sprintf("%v", v)
		case bool:
			secretResult[secretKey] = fmt.Sprintf("%t", v)
		default:
			return nil, fmt.Errorf("Key %s in secret for %s is not a string", secretKey, secretId)
		}
	}
	return secretResult, nil
}
