package main

import (
	"testing"
)

// TestIsPunct tests the punctuation detection function
func TestIsPunct(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		expected bool
	}{
		// Punctuation characters
		{"exclamation", '!', true},
		{"comma", ',', true},
		{"period", '.', true},
		{"question", '?', true},
		{"colon", ':', true},
		{"semicolon", ';', true},
		{"open paren", '(', true},
		{"close paren", ')', true},

		// Non-punctuation characters
		{"letter a", 'a', false},
		{"letter z", 'z', false},
		{"space", ' ', false},
		{"tab", '\t', false},
		{"newline", '\n', false},
		{"quote", '\'', false},
		{"digit 1", '1', false},
		{"digit 9", '9', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isPunct(tt.input)
			if result != tt.expected {
				t.Errorf("isPunct(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestIsSeparator tests the separator detection function
func TestIsSeparator(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		expected bool
	}{
		// Separators
		{"space", ' ', true},
		{"tab", '\t', true},
		{"newline", '\n', true},
		{"quote", '\'', true},
		{"comma", ',', true},
		{"period", '.', true},
		{"exclamation", '!', true},

		// Non-separators
		{"letter a", 'a', false},
		{"letter z", 'z', false},
		{"digit 1", '1', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSeparator(tt.input)
			if result != tt.expected {
				t.Errorf("isSeparator(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestIsVowelOrH tests the vowel/h detection function
func TestIsVowelOrH(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// Vowels
		{"lowercase a", "a", true},
		{"lowercase e", "e", true},
		{"lowercase i", "i", true},
		{"lowercase o", "o", true},
		{"lowercase u", "u", true},
		{"uppercase A", "A", false}, // Note: function expects lowercase

		// H sound
		{"h sound", "h", true},

		// Consonants
		{"b", "b", false},
		{"c", "c", false},
		{"z", "z", false},

		// Edge cases
		{"empty string", "", false},
		{"multiple chars", "ab", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isVowelOrH(tt.input)
			if result != tt.expected {
				t.Errorf("isVowelOrH(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestTransformWord tests word transformations
func TestTransformWord(t *testing.T) {
	tests := []struct {
		name     string
		word     string
		cmd      string
		expected string
	}{
		// Up command
		{"up: lowercase to uppercase", "hello", "up", "HELLO"},
		{"up: already uppercase", "HELLO", "up", "HELLO"},
		{"up: mixed case", "Hello", "up", "HELLO"},
		{"up: empty string", "", "up", ""},

		// Low command
		{"low: uppercase to lowercase", "WORLD", "low", "world"},
		{"low: already lowercase", "world", "low", "world"},
		{"low: mixed case", "WoRlD", "low", "world"},
		{"low: empty string", "", "low", ""},

		// Cap command
		{"cap: lowercase", "hello", "cap", "Hello"},
		{"cap: uppercase", "HELLO", "cap", "Hello"},
		{"cap: mixed case", "hElLo", "cap", "Hello"},
		{"cap: single character", "a", "cap", "A"},
		{"cap: empty string", "", "cap", ""},
		{"cap: with spaces", "hello world", "cap", "Hello world"},

		// Hex command
		{"hex: uppercase", "1A", "hex", "26"},
		{"hex: lowercase", "1a", "hex", "26"},
		{"hex: mixed case", "1A", "hex", "26"},
		{"hex: invalid", "XYZ", "hex", "XYZ"},
		{"hex: empty", "", "hex", ""},

		// Bin command
		{"bin: valid", "1010", "bin", "10"},
		{"bin: another valid", "1111", "bin", "15"},
		{"bin: invalid", "1234", "bin", "1234"},
		{"bin: empty", "", "bin", ""},

		// Unknown command
		{"unknown command", "hello", "unknown", "hello"},
		{"unknown marker", "world", "xyz", "world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := transformWord(tt.word, tt.cmd)
			if result != tt.expected {
				t.Errorf("transformWord(%q, %q) = %q, want %q",
					tt.word, tt.cmd, result, tt.expected)
			}
		})
	}
}

// TestTokenizer tests the tokenization function
func TestTokenizer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "simple words",
			input: "hello world",
			expected: []Token{
				{Type: "word", Value: "hello"},
				{Type: "word", Value: "world"},
			},
		},
		{
			name:  "multiple spaces",
			input: "hello    world",
			expected: []Token{
				{Type: "word", Value: "hello"},
				{Type: "word", Value: "world"},
			},
		},
		{
			name:  "punctuation",
			input: "hello, world!",
			expected: []Token{
				{Type: "word", Value: "hello"},
				{Type: "punct", Value: ","},
				{Type: "word", Value: "world"},
				{Type: "punct", Value: "!"},
			},
		},
		{
			name:  "multiple punctuation",
			input: "hello... world?!",
			expected: []Token{
				{Type: "word", Value: "hello"},
				{Type: "punct", Value: "."},
				{Type: "punct", Value: "."},
				{Type: "punct", Value: "."},
				{Type: "word", Value: "world"},
				{Type: "punct", Value: "?"},
				{Type: "punct", Value: "!"},
			},
		},
		{
			name:  "marker without number",
			input: "(up) hello",
			expected: []Token{
				{Type: "marker", Value: "(up)"},
				{Type: "word", Value: "hello"},
			},
		},
		{
			name:  "marker with number",
			input: "(up,2) hello world",
			expected: []Token{
				{Type: "marker", Value: "(up,2)"},
				{Type: "word", Value: "hello"},
				{Type: "word", Value: "world"},
			},
		},
		{
			name:  "multiple markers",
			input: "(up) hello (low) WORLD",
			expected: []Token{
				{Type: "marker", Value: "(up)"},
				{Type: "word", Value: "hello"},
				{Type: "marker", Value: "(low)"},
				{Type: "word", Value: "WORLD"},
			},
		},
		{
			name:  "quotes",
			input: "'hello' world",
			expected: []Token{
				{Type: "quote", Value: "'"},
				{Type: "word", Value: "hello"},
				{Type: "quote", Value: "'"},
				{Type: "word", Value: "world"},
			},
		},
		{
			name:  "newlines",
			input: "hello\nworld",
			expected: []Token{
				{Type: "word", Value: "hello"},
				{Type: "newline", Value: "\n"},
				{Type: "word", Value: "world"},
			},
		},
		{
			name:  "multiple newlines",
			input: "hello\n\nworld",
			expected: []Token{
				{Type: "word", Value: "hello"},
				{Type: "newline", Value: "\n"},
				{Type: "newline", Value: "\n"},
				{Type: "word", Value: "world"},
			},
		},
		{
			name:  "mixed content",
			input: "Hello (up,2) beautiful world! 'test'\nnew line",
			expected: []Token{
				{Type: "word", Value: "Hello"},
				{Type: "marker", Value: "(up,2)"},
				{Type: "word", Value: "beautiful"},
				{Type: "word", Value: "world"},
				{Type: "punct", Value: "!"},
				{Type: "quote", Value: "'"},
				{Type: "word", Value: "test"},
				{Type: "quote", Value: "'"},
				{Type: "newline", Value: "\n"},
				{Type: "word", Value: "new"},
				{Type: "word", Value: "line"},
			},
		},
		{
			name:     "empty input",
			input:    "",
			expected: []Token{},
		},
		{
			name:     "only spaces",
			input:    "   \t   ",
			expected: []Token{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Tokenizer(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("Tokenizer(%q) returned %d tokens, want %d",
					tt.input, len(result), len(tt.expected))
				// Print actual tokens for debugging
				t.Logf("Got tokens:")
				for i, tok := range result {
					t.Logf("  [%d] Type=%q, Value=%q", i, tok.Type, tok.Value)
				}
				return
			}

			for i := 0; i < len(result); i++ {
				if result[i].Type != tt.expected[i].Type {
					t.Errorf("Token %d: Type = %q, want %q",
						i, result[i].Type, tt.expected[i].Type)
				}
				if result[i].Value != tt.expected[i].Value {
					t.Errorf("Token %d: Value = %q, want %q",
						i, result[i].Value, tt.expected[i].Value)
				}
			}
		})
	}
}

// TestApplyMarkers tests the marker application function
func TestApplyMarkers(t *testing.T) {
	tests := []struct {
		name     string
		input    []Token
		expected []Token
	}{
		{
			name: "up marker single word",
			input: []Token{
				{Type: "word", Value: "hello"},
				{Type: "marker", Value: "(up)"},
			},
			expected: []Token{
				{Type: "word", Value: "HELLO"},
			},
		},
		{
			name: "up marker with count",
			input: []Token{
				{Type: "word", Value: "hello"},
				{Type: "word", Value: "beautiful"},
				{Type: "word", Value: "world"},
				{Type: "marker", Value: "(up,2)"},
			},
			expected: []Token{
				{Type: "word", Value: "hello"},
				{Type: "word", Value: "BEAUTIFUL"},
				{Type: "word", Value: "WORLD"},
			},
		},
		{
			name: "low marker",
			input: []Token{
				{Type: "word", Value: "WORLD"},
				{Type: "marker", Value: "(low)"},
			},
			expected: []Token{
				{Type: "word", Value: "world"},
			},
		},
		{
			name: "cap marker",
			input: []Token{
				{Type: "word", Value: "hello"},
				{Type: "marker", Value: "(cap)"},
			},
			expected: []Token{
				{Type: "word", Value: "Hello"},
			},
		},
		{
			name: "hex marker",
			input: []Token{
				{Type: "word", Value: "1A"},
				{Type: "marker", Value: "(hex)"},
			},
			expected: []Token{
				{Type: "word", Value: "26"},
			},
		},
		{
			name: "bin marker",
			input: []Token{
				{Type: "word", Value: "1010"},
				{Type: "marker", Value: "(bin)"},
			},
			expected: []Token{
				{Type: "word", Value: "10"},
			},
		},
		{
			name: "multiple markers",
			input: []Token{
				{Type: "word", Value: "hello"},
				{Type: "marker", Value: "(up)"},
				{Type: "word", Value: "world"},
				{Type: "marker", Value: "(low)"},
			},
			expected: []Token{
				{Type: "word", Value: "HELLO"},
				{Type: "word", Value: "world"},
			},
		},
		{
			name: "skip non-word tokens",
			input: []Token{
				{Type: "word", Value: "hello"},
				{Type: "punct", Value: ","},
				{Type: "word", Value: "world"},
				{Type: "marker", Value: "(up)"},
			},
			expected: []Token{
				{Type: "word", Value: "hello"},
				{Type: "punct", Value: ","},
				{Type: "word", Value: "WORLD"},
			},
		},
		{
			name: "not enough words",
			input: []Token{
				{Type: "word", Value: "hello"},
				{Type: "marker", Value: "(up,5)"},
			},
			expected: []Token{
				{Type: "word", Value: "HELLO"},
			},
		},
		{
			name: "no markers",
			input: []Token{
				{Type: "word", Value: "hello"},
				{Type: "word", Value: "world"},
			},
			expected: []Token{
				{Type: "word", Value: "hello"},
				{Type: "word", Value: "world"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := applyMarkers(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d tokens, got %d", len(tt.expected), len(result))
				return
			}

			for i := 0; i < len(result); i++ {
				if result[i].Type != tt.expected[i].Type {
					t.Errorf("Token %d: Type = %q, want %q",
						i, result[i].Type, tt.expected[i].Type)
				}
				if result[i].Value != tt.expected[i].Value {
					t.Errorf("Token %d: Value = %q, want %q",
						i, result[i].Value, tt.expected[i].Value)
				}
			}
		})
	}
}

// TestFixArticles tests the article transformation function
func TestFixArticles(t *testing.T) {
	tests := []struct {
		name     string
		input    []Token
		expected []Token
	}{
		{
			name: "a to an before vowel",
			input: []Token{
				{Type: "word", Value: "a"},
				{Type: "word", Value: "apple"},
			},
			expected: []Token{
				{Type: "word", Value: "an"},
				{Type: "word", Value: "apple"},
			},
		},
		{
			name: "A to An before vowel",
			input: []Token{
				{Type: "word", Value: "A"},
				{Type: "word", Value: "apple"},
			},
			expected: []Token{
				{Type: "word", Value: "An"},
				{Type: "word", Value: "apple"},
			},
		},
		{
			name: "a before consonant - no change",
			input: []Token{
				{Type: "word", Value: "a"},
				{Type: "word", Value: "book"},
			},
			expected: []Token{
				{Type: "word", Value: "a"},
				{Type: "word", Value: "book"},
			},
		},
		{
			name: "skip punctuation",
			input: []Token{
				{Type: "word", Value: "a"},
				{Type: "punct", Value: ","},
				{Type: "word", Value: "apple"},
			},
			expected: []Token{
				{Type: "word", Value: "an"},
				{Type: "punct", Value: ","},
				{Type: "word", Value: "apple"},
			},
		},
		{
			name: "skip quotes",
			input: []Token{
				{Type: "word", Value: "a"},
				{Type: "quote", Value: "'"},
				{Type: "word", Value: "apple"},
			},
			expected: []Token{
				{Type: "word", Value: "an"},
				{Type: "quote", Value: "'"},
				{Type: "word", Value: "apple"},
			},
		},
		{
			name: "h sound",
			input: []Token{
				{Type: "word", Value: "a"},
				{Type: "word", Value: "hour"},
			},
			expected: []Token{
				{Type: "word", Value: "an"},
				{Type: "word", Value: "hour"},
			},
		},
		{
			name: "multiple articles",
			input: []Token{
				{Type: "word", Value: "a"},
				{Type: "word", Value: "book"},
				{Type: "word", Value: "a"},
				{Type: "word", Value: "apple"},
			},
			expected: []Token{
				{Type: "word", Value: "a"},
				{Type: "word", Value: "book"},
				{Type: "word", Value: "an"},
				{Type: "word", Value: "apple"},
			},
		},
		{
			name: "a at end - no change",
			input: []Token{
				{Type: "word", Value: "a"},
			},
			expected: []Token{
				{Type: "word", Value: "a"},
			},
		},
		{
			name: "skip markers",
			input: []Token{
				{Type: "word", Value: "a"},
				{Type: "marker", Value: "(up)"},
				{Type: "word", Value: "apple"},
			},
			expected: []Token{
				{Type: "word", Value: "an"},
				{Type: "marker", Value: "(up)"},
				{Type: "word", Value: "apple"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fixArticles(tt.input)

			for i := 0; i < len(result); i++ {
				if result[i].Value != tt.expected[i].Value {
					t.Errorf("Token %d: Value = %q, want %q",
						i, result[i].Value, tt.expected[i].Value)
				}
			}
		})
	}
}

// TestAssemble tests the assembly function
func TestAssemble(t *testing.T) {
	tests := []struct {
		name     string
		input    []Token
		expected string
	}{
		{
			name: "simple words",
			input: []Token{
				{Type: "word", Value: "hello"},
				{Type: "word", Value: "world"},
			},
			expected: "hello world",
		},
		{
			name: "punctuation",
			input: []Token{
				{Type: "word", Value: "hello"},
				{Type: "punct", Value: ","},
				{Type: "word", Value: "world"},
			},
			expected: "hello, world",
		},
		{
			name: "multiple punctuation",
			input: []Token{
				{Type: "word", Value: "hello"},
				{Type: "punct", Value: "."},
				{Type: "punct", Value: "."},
				{Type: "punct", Value: "."},
				{Type: "word", Value: "world"},
			},
			expected: "hello... world",
		},
		{
			name: "quotes - simple",
			input: []Token{
				{Type: "quote", Value: "'"},
				{Type: "word", Value: "hello"},
				{Type: "quote", Value: "'"},
			},
			expected: "'hello'",
		},
		{
			name: "quotes with space before",
			input: []Token{
				{Type: "word", Value: "say"},
				{Type: "quote", Value: "'"},
				{Type: "word", Value: "hello"},
				{Type: "quote", Value: "'"},
			},
			expected: "say 'hello'",
		},
		{
			name: "newlines",
			input: []Token{
				{Type: "word", Value: "hello"},
				{Type: "newline", Value: "\n"},
				{Type: "word", Value: "world"},
			},
			expected: "hello\nworld",
		},
		{
			name: "multiple newlines",
			input: []Token{
				{Type: "word", Value: "hello"},
				{Type: "newline", Value: "\n"},
				{Type: "newline", Value: "\n"},
				{Type: "word", Value: "world"},
			},
			expected: "hello\n\nworld",
		},
		{
			name: "punctuation after quote",
			input: []Token{
				{Type: "word", Value: "say"},
				{Type: "quote", Value: "'"},
				{Type: "word", Value: "hello"},
				{Type: "quote", Value: "'"},
				{Type: "punct", Value: ","},
				{Type: "word", Value: "world"},
			},
			expected: "say 'hello', world",
		},
		{
			name: "mixed content",
			input: []Token{
				{Type: "word", Value: "Hello"},
				{Type: "punct", Value: "!"},
				{Type: "word", Value: "How"},
				{Type: "quote", Value: "'"},
				{Type: "word", Value: "are"},
				{Type: "quote", Value: "'"},
				{Type: "word", Value: "you"},
				{Type: "punct", Value: "?"},
			},
			expected: "Hello! How 'are' you?",
		},
		{
			name:     "empty input",
			input:    []Token{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := assemble(tt.input)
			if result != tt.expected {
				t.Errorf("assemble() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// TestIntegration tests the complete pipeline
func TestIntegration(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic marker",
			input:    "hello (up) world",
			expected: "HELLO world",
		},
		{
			name:     "marker with count",
			input:    "hello beautiful world (up,2)",
			expected: "hello BEAUTIFUL WORLD",
		},
		{
			name:     "article conversion",
			input:    "a apple",
			expected: "an apple",
		},
		{
			name:     "article with punctuation",
			input:    "a, apple",
			expected: "an, apple",
		},
		{
			name:     "hex conversion",
			input:    "(hex) 1A",
			expected: "26",
		},
		{
			name:     "bin conversion",
			input:    "(bin) 1010",
			expected: "10",
		},
		{
			name:     "cap conversion",
			input:    "(cap) hello",
			expected: "Hello",
		},
		{
			name:     "complex example",
			input:    "hello (up,2) beautiful world - a apple (hex) 1A",
			expected: "hello BEAUTIFUL WORLD - an apple 26",
		},
		{
			name:     "quotes with markers",
			input:    "'hello' (up) world",
			expected: "'hello' WORLD",
		},
		{
			name:     "multiple transformations",
			input:    "a (cap) beautiful (hex) 1A (bin) 1010 world",
			expected: "a Beautiful 26 10 world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := Tokenizer(tt.input)
			tokens = applyMarkers(tokens)
			tokens = fixArticles(tokens)
			result := assemble(tokens)

			if result != tt.expected {
				t.Errorf("Full pipeline(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// Benchmark tests for performance
func BenchmarkTokenizer(b *testing.B) {
	input := "hello (up,2) beautiful world! 'test' (hex) 1A (bin) 1010"
	for i := 0; i < b.N; i++ {
		Tokenizer(input)
	}
}

func BenchmarkApplyMarkers(b *testing.B) {
	tokens := Tokenizer("hello (up,2) beautiful world! 'test' (hex) 1A (bin) 1010")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		applyMarkers(tokens)
	}
}

func BenchmarkFullPipeline(b *testing.B) {
	input := "hello (up,2) beautiful world! 'test' (hex) 1A (bin) 1010 a apple"
	for i := 0; i < b.N; i++ {
		tokens := Tokenizer(input)
		tokens = applyMarkers(tokens)
		tokens = fixArticles(tokens)
		_ = assemble(tokens)
	}
}
