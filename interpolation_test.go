package secretsenv

import (
	"testing"
)

func TestInterpolator_Interpolate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		vars    map[string]string
		want    string
		wantErr bool
	}{
		{
			name:    "Simple interpolation",
			input:   "Hello, ${name}!",
			vars:    map[string]string{"name": "Alice"},
			want:    "Hello, Alice!",
			wantErr: false,
		},
		{
			name:    "Escaped variable",
			input:   "This is a \\${escaped} variable",
			vars:    map[string]string{"escaped": "test"},
			want:    "This is a ${escaped} variable",
			wantErr: false,
		},
		{
			name:    "Multiple variables",
			input:   "${greet}, ${name}!",
			vars:    map[string]string{"greet": "Hello", "name": "Bob"},
			want:    "Hello, Bob!",
			wantErr: false,
		},
		{
			name:    "Missing variable",
			input:   "Hello, ${name}!",
			vars:    map[string]string{},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Mixed variables and plain text",
			input:   "Value: $value, Text: plain text",
			vars:    map[string]string{"value": "123"},
			want:    "Value: 123, Text: plain text",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Interpolator{}
			got, err := i.Interpolate(tt.input, tt.vars)
			if err != nil != tt.wantErr {
				t.Errorf("Interpolator.Interpolate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Interpolator.Interpolate() = %v, want %v", got, tt.want)
			}
		})
	}
}
