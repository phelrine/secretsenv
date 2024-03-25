package secretsenv

import (
	"errors"
	"regexp"
	"strings"
)

type Interpolator struct {
}

// ErrVariableNotFound is the error returned when a variable is not found
//
// It is returned by the Interpolate method when a variable is not found in the vars map.
var ErrVariableNotFound = errors.New("variable not found")

// Interpolate replaces variables in the input string with the values from the vars map
//
// It uses the vars map to replace variables in the input string. The variables are
// specified using the ${VAR} or $VAR syntax. If a variable is not found in the vars
// map, it returns an error.
// It also supports escaping the dollar sign with a backslash. For example, \$VAR
// will be replaced with $VAR.
func (i *Interpolator) Interpolate(input string, vars map[string]string) (string, error) {
	// regex to match ${VAR} or $VAR
	re := regexp.MustCompile(`(\\)?\$\{?[a-zA-Z0-9_]+}?`)
	var err error
	// replace all matches with the value of the variable
	result := re.ReplaceAllStringFunc(input, func(match string) string {
		if strings.HasPrefix(match, "\\") {
			return match[1:]
		}
		varName := strings.TrimRight(strings.TrimLeft(match, "${"), "}")
		value, ok := vars[varName]
		if ok {
			return value
		}
		if err == nil {
			err = ErrVariableNotFound
		}
		err = errors.Join(err, errors.New("variable not found: "+varName))
		return ""
	})
	if err != nil {
		return "", err
	}
	// replace escaped dollar signs
	result = strings.ReplaceAll(result, "\\$", "$")
	return result, nil
}
