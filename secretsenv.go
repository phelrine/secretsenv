package secretsenv

import "fmt"

type SecretsEnv struct {
	Loaders map[string]SecretLoader
}

func (s *SecretsEnv) Load(option SecretOption) (SecretMapping, error) {
	loader, ok := s.Loaders[option.Type]
	if !ok {
		return nil, fmt.Errorf(
			"specified loader %s not found",
			option.Type,
		)
	}
	secrets, err := loader.Load(option.SecretId, option.Option)
	if err != nil {
		return nil, err
	}
	// map the secret values to the environment variables
	mappingResult, err := s.mapSecrets(option, secrets)
	if err != nil {
		return nil, err
	}
	return mappingResult, nil
}

// mapSecrets maps the secret values to the environment variables
//
// It uses the mapping defined in the configuration file to map the secret values
// to the environment variables. If a mapping is not defined, it uses the variable
// name as the key.
func (s *SecretsEnv) mapSecrets(item SecretOption, secrets SecretResult) (SecretMapping, error) {
	mappingResult := make(SecretMapping)
	interpolator := &Interpolator{}
	for variableName, secretKey := range item.Mapping {
		if secretKey == nil {
			value, ok := secrets[variableName]
			if !ok {
				return nil, fmt.Errorf(
					"specified key %s not found in the secret %s",
					variableName,
					item.SecretId,
				)
			}
			mappingResult[variableName] = value
			continue
		}
		value, err := interpolator.Interpolate(*secretKey, secrets)
		if err != nil {
			return nil, err
		}
		mappingResult[variableName] = value
	}
	return mappingResult, nil
}
