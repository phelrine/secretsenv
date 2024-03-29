package secretsenv

import (
	"fmt"
	"reflect"
	"testing"
)

// MockLoader implements the SecretLoader interface for testing purposes.
type MockLoader struct {
	ShouldFail  bool
	ReturnError error
	Secrets     SecretResult
}

func (m *MockLoader) Load(secretId string, options LoaderOption) (SecretResult, error) {
	if m.ShouldFail {
		return nil, m.ReturnError
	}
	return m.Secrets, nil
}

func TestSecretsEnv_Load(t *testing.T) {
	val := "$SecretField"
	tests := []struct {
		name           string
		option         SecretOption
		mockLoader     *MockLoader
		expectedResult SecretMapping
		expectedError  string
	}{
		{
			name: "Loader not found",
			option: SecretOption{
				Type: "unknown",
			},
			mockLoader:     nil, // Loader not registered.
			expectedResult: nil,
			expectedError:  "specified loader unknown not found",
		},
		{
			name: "Loader Load error",
			option: SecretOption{
				Type: "mock",
			},
			mockLoader: &MockLoader{
				ShouldFail:  true,
				ReturnError: fmt.Errorf("failed to load secret"),
			},
			expectedResult: nil,
			expectedError:  "failed to load secret",
		},
		{
			name: "Variable not found (load results does not contain the specified key)",
			option: SecretOption{
				Type:     "mock",
				SecretId: "secretId",
				Mapping: map[string]*string{
					"EnvName": &val,
				},
			},
			mockLoader: &MockLoader{
				ShouldFail: false,
				Secrets:    SecretResult{"SecretFieldNG": "SecretValue"},
			},
			expectedResult: nil,
			expectedError:  "variable not found: SecretField",
		},
		{
			name: "Variable not found (load results does not contain the same key as the specified key)",
			option: SecretOption{
				Type:     "mock",
				SecretId: "secretId",
				Mapping: map[string]*string{
					"EnvName": nil,
				},
			},
			mockLoader: &MockLoader{
				ShouldFail: false,
				Secrets:    SecretResult{"SecretField": "SecretValue"},
			},
			expectedResult: nil,
			expectedError:  "specified key EnvName not found in the secret secretId",
		},
		{
			name: "Successful load and map",
			option: SecretOption{
				Type: "mock",
				Mapping: map[string]*string{
					"EnvName":  &val,
					"EnvName2": nil,
				},
			},
			mockLoader: &MockLoader{
				ShouldFail: false,
				Secrets: SecretResult{
					"SecretField": "SecretValue",
					"EnvName2":    "SecretValue2",
				},
			},
			expectedResult: SecretMapping{
				"EnvName":  "SecretValue",
				"EnvName2": "SecretValue2",
			},
			expectedError: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := SecretsEnv{
				Loaders: map[string]SecretLoader{},
			}
			if test.mockLoader != nil {
				s.Loaders["mock"] = test.mockLoader
			}
			result, err := s.Load(test.option)
			if err != nil {
				if uw, ok := err.(interface{ Unwrap() []error }); ok {
					errs := uw.Unwrap()
					if errs[len(errs)-1].Error() != test.expectedError {
						t.Errorf("Expected error: %s, got: %s", test.expectedError, err)
					}
				} else if err.Error() != test.expectedError {
					t.Errorf("Expected error: %s, got: %s", test.expectedError, err)
				}
			} else if err == nil && test.expectedError != "" {
				t.Errorf("Expected error: %s, but no error was received", test.expectedError)
			}
			if !reflect.DeepEqual(result, test.expectedResult) {
				t.Errorf("Expected result: %v, got: %v", test.expectedResult, result)
			}
		})
	}
}
