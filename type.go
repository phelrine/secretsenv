package secretsenv

// SecretOption represents the options for the secret loader
//
// It is a map of options where the key is the name of the option
// and the value is the value of the option.
// The options are specific to the secret loader.
// For example, the AWS secret loader may have an option for the region.
// The options are passed to the secret loader when loading the secret.
// The secret loader may use the options to configure the behavior of the loader.
type SecretOption map[string]interface{}

// SecretMapping represents the mapping of the secret values to the environment variables
//
// It is a map of environment variables where the key is the name of the environment variable
// and the value is the name of the secret value.
type SecretMapping map[string]string

// SecretResult represents the result of loading a secret
//
// It is a map of secret values where the key is the name of the secret value
// and the value is the value of the secret value.
type SecretResult map[string]string

// SecretLoader represents a secret loader
//
// It is an interface that defines the Load method, which loads a secret
// from a secret store and returns the result.
type SecretLoader interface {
	Load(secretId string, option SecretOption) (SecretResult, error)
}
