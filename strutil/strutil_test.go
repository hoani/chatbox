package strutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplitWidth(t *testing.T) {
	testCases := []struct {
		name     string
		in       string
		width    int
		expected []string
	}{
		{
			name:     "empty string",
			in:       "",
			width:    10,
			expected: []string{},
		},
		{
			name:     "short string",
			in:       "hello",
			width:    10,
			expected: []string{"hello"},
		},
		{
			name:     "long string",
			in:       "hello world",
			width:    5,
			expected: []string{"hello", "world"},
		},
		{
			name:     "long string with spaces",
			in:       "hello world, how are you today?",
			width:    7,
			expected: []string{"hello", "world,", "how are", "you", "today?"},
		},
		{
			name:     "long string with spaces and tabs",
			in:       "hello world,\thow are you today?",
			width:    7,
			expected: []string{"hello", "world,", "how are", "you", "today?"},
		},
		{
			name:     "words longer than max width",
			in:       "hello hippopotamus, how are you today?",
			width:    10,
			expected: []string{"hello", "hippopotam", "us, how", "are you", "today?"},
		},
		{
			name:     "really long words",
			in:       "hippopotamus girraffe elephant",
			width:    4,
			expected: []string{"hipp", "opot", "amus", "girr", "affe", "elep", "hant"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := SplitWidth(tc.in, tc.width)
			require.Equal(t, tc.expected, actual)
		})
	}
}
