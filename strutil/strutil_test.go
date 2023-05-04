package strutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSimplify(t *testing.T) {
	testCases := []struct {
		name     string
		in       string
		expected string
	}{
		{
			name:     "empty string",
			in:       "",
			expected: "",
		},
		{
			name:     "basic sentence",
			in:       "hello world",
			expected: "hello world",
		},
		{
			name:     "additional letters",
			in:       "hhhhello worlllddd",
			expected: "hello world",
		},
		{
			name:     "limit vowels to 2",
			in:       "hhhheeeelloo worllld",
			expected: "heelloo world",
		},
		{
			name:     "limit punctuation to 1",
			in:       "hello world!!!!!!",
			expected: "hello world!",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, Simplify(tc.in))
		})
	}
}

func TestSplitBrackets(t *testing.T) {
	testCases := []struct {
		name     string
		in       string
		expected []string
	}{
		{
			name:     "empty string",
			in:       "",
			expected: []string{},
		},
		{
			name:     "no brackets",
			in:       "hello world",
			expected: []string{"hello world"},
		},
		{
			name:     "curved bracket complete",
			in:       "hello (world) hello",
			expected: []string{"hello", "(world)", "hello"},
		},
		{
			name:     "curved bracket incomplete",
			in:       "hello (world",
			expected: []string{"hello", "(world"},
		},
		{
			name:     "curved terminating bracket",
			in:       "hello) world",
			expected: []string{"hello) world"},
		},
		{
			name:     "curly bracket complete",
			in:       "hello {world}",
			expected: []string{"hello", "{world}"},
		},
		{
			name:     "square bracket complete",
			in:       "hello [world]",
			expected: []string{"hello", "[world]"},
		},
		{
			name:     "angle bracket complete",
			in:       "hello <world>",
			expected: []string{"hello", "<world>"},
		},
		{
			name:     "curved and curly bracket complete",
			in:       "hello (world) {how *are* you today?}",
			expected: []string{"hello", "(world)", "{how *are* you today?}"},
		},
		{
			name:     "nested brackets",
			in:       "hello (world {how (are) you today?})",
			expected: []string{"hello", "(world {how (are) you today?})"},
		},
		{
			name:     "asterixes",
			in:       "hello *world* how are you today?",
			expected: []string{"hello", "*world*", "how are you today?"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := SplitBrackets(tc.in)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestSplitSentences(t *testing.T) {
	testCases := []struct {
		name     string
		in       string
		expected []string
	}{
		{
			name:     "empty string",
			in:       "",
			expected: []string{},
		},
		{
			name:     "no sentences",
			in:       "hello world",
			expected: []string{"hello world"},
		},
		{
			name:     "one sentence",
			in:       "hello world.",
			expected: []string{"hello world."},
		},
		{
			name:     "two sentences",
			in:       "hello world. how are you today?",
			expected: []string{"hello world.", "how are you today?"},
		},
		{
			name:     "many sentences with punctuation",
			in:       "How are you today? Leave me alone! I'm sorry.",
			expected: []string{"How are you today?", "Leave me alone!", "I'm sorry."},
		},
		{
			name:     "groups punctuation",
			in:       "How are you today??? Leave me alone!!!!! I'm sorry...",
			expected: []string{"How are you today???", "Leave me alone!!!!!", "I'm sorry..."},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := SplitSentences(tc.in)
			require.Equal(t, tc.expected, actual)
		})
	}
}

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
