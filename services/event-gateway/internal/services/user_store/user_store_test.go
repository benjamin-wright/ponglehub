package user_store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListIds(t *testing.T) {
	for _, test := range []struct {
		name       string
		nameLookup map[string]string
		id         string
		expected   []string
	}{
		{
			name:       "empty",
			nameLookup: map[string]string{},
			id:         "abc123",
			expected:   []string{},
		},
		{
			name:       "only self",
			nameLookup: map[string]string{"abc123": "username"},
			id:         "abc123",
			expected:   []string{},
		},
		{
			name: "missing self",
			nameLookup: map[string]string{
				"def456": "other",
				"ghi789": "users",
			},
			id:       "abc123",
			expected: []string{"def456", "ghi789"},
		},
		{
			name: "filters self",
			nameLookup: map[string]string{
				"abc123": "username",
				"def456": "other",
				"ghi789": "users",
			},
			id:       "abc123",
			expected: []string{"def456", "ghi789"},
		},
	} {
		t.Run(test.name, func(u *testing.T) {
			store := &Store{nameLookup: test.nameLookup}
			actual := store.ListIDs(test.id)

			assert.Equal(u, test.expected, actual)
		})
	}
}
