package user_store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFriends(t *testing.T) {
	for _, test := range []struct {
		name       string
		nameLookup map[string]string
		id         string
		expected   map[string]string
	}{
		{
			name:       "empty",
			nameLookup: map[string]string{},
			id:         "abc123",
			expected:   map[string]string{},
		},
		{
			name:       "only self",
			nameLookup: map[string]string{"abc123": "username"},
			id:         "abc123",
			expected:   map[string]string{},
		},
		{
			name: "missing self",
			nameLookup: map[string]string{
				"def456": "other",
				"ghi789": "users",
			},
			id:       "abc123",
			expected: map[string]string{"def456": "other", "ghi789": "users"},
		},
		{
			name: "filters self",
			nameLookup: map[string]string{
				"abc123": "username",
				"def456": "other",
				"ghi789": "users",
			},
			id:       "abc123",
			expected: map[string]string{"def456": "other", "ghi789": "users"},
		},
	} {
		t.Run(test.name, func(u *testing.T) {
			store := &Store{nameLookup: test.nameLookup}
			actual := store.GetFriends(test.id)

			assert.Equal(u, test.expected, actual)
		})
	}
}
