package eval

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvaluate(t *testing.T) {
	type test struct {
		name    string
		source  string
		defines map[string]string
		want    bool
	}

	tests := []test{
		{
			name:    "empty",
			source:  "",
			defines: map[string]string{},
			want:    false,
		},
		{
			name:    "whitespace",
			source:  " ",
			defines: map[string]string{},
			want:    false,
		},
		{
			name:    "id undefined",
			source:  "abc",
			defines: map[string]string{},
			want:    false,
		},
		{
			name:    "id defined",
			source:  "abc",
			defines: map[string]string{"abc": "1"},
			want:    true,
		},
		{
			name:    "id double resolve",
			source:  "abc",
			defines: map[string]string{"abc": "def", "def": "1"},
			want:    true,
		},
		{
			name:    "id double resolve into undefined",
			source:  "abc",
			defines: map[string]string{"abc": "def", "def": "ghi"},
			want:    false,
		},
		{
			name:    "circular",
			source:  "abc",
			defines: map[string]string{"abc": "def", "def": "abc"},
			want:    false,
		},
	}

	for _, tt := range tests {
		actual := Evaluate(tt.source, tt.defines)
		assert.Equal(t, tt.want, actual, tt.name)
	}
}
