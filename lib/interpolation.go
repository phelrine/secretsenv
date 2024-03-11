package lib

import (
	"errors"
	"regexp"
	"strings"
)

type Interpolator struct {
}

var ErrVariableNotFound = errors.New("variable not found")

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
