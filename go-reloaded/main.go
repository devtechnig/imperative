package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type Token struct {
	Type  string // "word", "punct", "quote", "newline", "marker"
	Value string
}

func Tokenizer(text string) []Token {
	var tokens []Token
	i := 0
	n := len(text)

	for i < n {
		ch := text[i]
		if ch == '(' {
			end := strings.IndexByte(text[i:], ')')
			if end != -1 {
				marker := text[i : i+end+1]
				content := text[i+1 : i+end]
				parts := strings.SplitN(content, ",", 2)
				markerWord := strings.TrimSpace(parts[0])
				switch markerWord {
				case "up", "low", "cap", "hex", "bin":
					tokens = append(tokens, Token{
						Type:  "marker",
						Value: marker,
					})
					i = i + end + 1
					continue
				}
			}
		}

		switch {
		case ch == '\'':
			tokens = append(tokens, Token{
				Type:  "quote",
				Value: "'",
			})
			i++
		case ch == '\n':
			tokens = append(tokens, Token{
				Type:  "newline",
				Value: "\n",
			})
			i++
		case ch == ' ' || ch == '\t':
			i++
		case isPunct(ch):
			tokens = append(tokens, Token{
				Type:  "punct",
				Value: string(ch),
			})
			i++
		default:
			start := i
			for i < n && !isSeparator(text[i]) {
				i++
			}
			word := text[start:i]
			tokens = append(tokens, Token{
				Type:  "word",
				Value: word,
			})
		}
	}
	return tokens
}

func isPunct(c byte) bool {
	switch c {
	case '!', ',', '.', '?', ':', ';', '(', ')':
		return true
	}
	return false
}

func isSeparator(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\'' || isPunct(c)
}

func applyMarkers(tokens []Token) []Token {
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Type != "marker" {
			continue
		}
		// parse marker
		inner := tokens[i].Value[1 : len(tokens[i].Value)-1]
		parts := strings.SplitN(inner, ",", 2)
		cmd := parts[0]
		count := 1
		if len(parts) > 1 {
			if c, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
				count = c
			}
		}

		// apply to previous words (backward only)
		applied := 0
		for j := i - 1; j >= 0 && applied < count; j-- {
			if tokens[j].Type == "word" {
				tokens[j].Value = transformWord(tokens[j].Value, cmd)
				applied++
			}
		}

		// remove marker
		tokens = append(tokens[:i], tokens[i+1:]...)
		i-- // adjust index after removal
	}
	return tokens
}

func transformWord(word, cmd string) string {
	switch cmd {
	case "up":
		return strings.ToUpper(word)
	case "low":
		return strings.ToLower(word)
	case "cap":
		if word == "" {
			return word
		}
		runes := []rune(word)
		res := string(unicode.ToUpper(runes[0]))
		if len(runes) > 1 {
			res += strings.ToLower(string(runes[1:]))
		}
		return res
	case "hex":
		if dec, err := strconv.ParseInt(word, 16, 64); err == nil {
			return strconv.FormatInt(dec, 10)
		}
		return word
	case "bin":
		if dec, err := strconv.ParseInt(word, 2, 64); err == nil {
			return strconv.FormatInt(dec, 10)
		}
		return word
	}
	return word
}

func fixArticles(tokens []Token) []Token {
	for i := range tokens {
		if tokens[i].Type != "word" {
			continue
		}
		val := strings.ToLower(tokens[i].Value)
		if val != "a" {
			continue
		}

		for j := i + 1; j < len(tokens); j++ {
			if tokens[j].Type != "word" {
				continue
			}

			nextWord := tokens[j].Value
			if len(nextWord) == 0 {
				break
			}

			firstChar := strings.ToLower(string(nextWord[0]))
			if isVowelOrH(firstChar) {
				switch tokens[i].Value {
				case "a":
					tokens[i].Value = "an"
				case "A":
					tokens[i].Value = "An"
				}
			}
			break
		}
	}
	return tokens
}

func isVowelOrH(ch string) bool {
	switch ch {
	case "a", "e", "i", "o", "u", "h":
		return true
	}
	return false
}

func assemble(tokens []Token) string {
	var builder strings.Builder
	quoteOpen := false

	for i, token := range tokens {
		switch token.Type {
		case "word":
			if i > 0 {
				prev := tokens[i-1]
				if prev.Type == "word" || (prev.Type == "quote" && !quoteOpen) {
					builder.WriteString(" ")
				}
			}
			builder.WriteString(token.Value)

		case "punct":
			builder.WriteString(token.Value)
			if i+1 < len(tokens) {
				next := tokens[i+1]
				// Amended: Prevent adding a space if next token is punctuation, newline, or a closing quote
				if next.Type != "punct" && next.Type != "newline" && !(next.Type == "quote" && quoteOpen) {
					builder.WriteString(" ")
				}
			}

		case "quote":
			if !quoteOpen {
				// Opening quote logic
				if i > 0 {
					prev := tokens[i-1]
					if prev.Type == "word" || (prev.Type == "quote" && !quoteOpen) {
						builder.WriteString(" ")
					}
				}
				builder.WriteString(token.Value)
				quoteOpen = true
			} else {
				// Closing quote logic
				builder.WriteString(token.Value)
				quoteOpen = false
			}

		case "newline":
			builder.WriteString("\n")
		}
	}
	return strings.TrimSpace(builder.String())
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %v <input file> <output file>\n", os.Args[0])
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	err := run(inputPath, outputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully processed %s -> %s\n", inputPath, outputPath)
}

func run(inputPath, outputPath string) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	text := string(data)
	tokens := Tokenizer(text)
	tokens = applyMarkers(tokens)
	tokens = fixArticles(tokens)
	result := assemble(tokens)

	return os.WriteFile(outputPath, []byte(result), 0644)

}
