package secretsenv

type LoaderOption map[string]interface{}

// SecretMapping represents the mapping of secret values to environment variables
//
// It is a map of environment variables where the key is the name of the
// environment variable and the value is the value of the environment variable.
type SecretMapping map[string]string

// SecretOption represents the configuration for loading a secret
//
// It is a struct that contains the type of the loader, the secret ID, the
// loader options, and the mapping of secret values to environment variables.
type SecretOption struct {
	Type     string             `yaml:"type,omitempty"`
	SecretId string             `yaml:"secretId"`
	Option   LoaderOption       `yaml:"option"`
	Mapping  map[string]*string `yaml:",inline"`
}

type SecretResult map[string]string

// SecretLoader represents a secret loader
//
// It is an interface that defines the Load method, which loads a secret
// from a secret store and returns the result.
type SecretLoader interface {
	Load(secretId string, option LoaderOption) (SecretResult, error)
}
