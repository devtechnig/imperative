package main

import (
	"testing"
)

// TestOfficialExamples tests the exact input/output pairs provided in the assignment.
func TestOfficialExamples(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "markers: cap, up, cap with count, low with count",
			input: "it (cap) was the best of times, it was the worst of times (up) , it was the age of wisdom, it was the age of foolishness (cap, 6) , it was the epoch of belief, it was the epoch of incredulity, it was the season of Light, it was the season of darkness, it was the spring of hope, IT WAS THE (low, 3) winter of despair.",
			want:  "It was the best of times, it was the worst of TIMES, it was the age of wisdom, It Was The Age Of Foolishness, it was the epoch of belief, it was the epoch of incredulity, it was the season of Light, it was the season of darkness, it was the spring of hope, it was the winter of despair.",
		},
		{
			name:  "hex and bin conversion",
			input: "Simply add 42 (hex) and 10 (bin) and you will see the result is 68.",
			want:  "Simply add 66 and 2 and you will see the result is 68.",
		},
		{
			name:  "a -> an before vowel sound",
			input: "There is no greater agony than bearing a untold story inside you.",
			want:  "There is no greater agony than bearing an untold story inside you.",
		},
		{
			name:  "punctuation spacing",
			input: "Punctuation tests are ... kinda boring ,what do you think ?",
			want:  "Punctuation tests are... kinda boring, what do you think?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := Tokenizer(tt.input)
			tokens = applyMarkers(tokens)
			tokens = fixArticles(tokens)
			got := assemble(tokens)

			if got != tt.want {
				t.Errorf("want:\n%q\ngot:\n%q", tt.want, got)
			}
		})
	}
}

// Additional test for a simple marker (backward only)
func TestBasicMarker(t *testing.T) {
	input := "hello (up) world"
	want := "HELLO world"
	tokens := Tokenizer(input)
	tokens = applyMarkers(tokens)
	got := assemble(tokens)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// Test that marker with count applies to previous words only
func TestMarkerWithCount(t *testing.T) {
	input := "hello beautiful world (up,2)"
	want := "hello BEAUTIFUL WORLD"
	tokens := Tokenizer(input)
	tokens = applyMarkers(tokens)
	got := assemble(tokens)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
