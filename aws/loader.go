package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/phelrine/secretsenv"
)

type SecretsManagerLoader struct {
}

// NewSecretsManagerLoader creates a new SecretsManagerLoader
//
// It creates a new SecretsManagerLoader with default values.
func NewSecretsManagerLoader() *SecretsManagerLoader {
	return &SecretsManagerLoader{}
}

// Load loads a secret from AWS Secrets Manager
//
// It loads a secret from AWS Secrets Manager using the specified secret ID.
// It returns the secret values as a map of strings where the key is the name
// of the secret value and the value is the value of the secret value.
// It returns an error if the secret cannot be loaded. The error may be due to
// the secret not existing or the user not having permission to access the secret.
func (l *SecretsManagerLoader) Load(secretId string, option secretsenv.SecretOption) (secretsenv.SecretResult, error) {
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
	secretResult := make(secretsenv.SecretResult)
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
