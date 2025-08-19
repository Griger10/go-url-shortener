package random

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRandomString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		size int
	}{
		{
			"size=5",
			5,
		},
		{
			"size=10",
			10,
		},
		{
			"size=20",
			20,
		},
		{
			"size=30",
			30,
		},
	}

	isAllowed := func(s string) bool {
		for _, r := range s {
			if !(r >= 'A' && r <= 'Z' ||
				r >= 'a' && r <= 'z' ||
				r >= '0' && r <= '9') {
				return false
			}
		}
		return true
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewRandomString(tt.size)
			assert.Len(t, s, tt.size)
			assert.True(t, isAllowed(s))

			const N = 200
			seen := make(map[string]struct{}, N)
			for i := 0; i < N; i++ {
				seen[NewRandomString(tt.size)] = struct{}{}
			}
			assert.Greater(t, len(seen), N*9/10)
		})
	}
}
