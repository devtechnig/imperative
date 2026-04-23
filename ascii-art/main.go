package main

import (
	"fmt"
	"os"
	"strings"
)

func parseBanner(lines []string) (map[byte][]string, error) {
	asciiMap := make(map[byte][]string)

	// Remove any empty lines at the beginning or end
	start := 0
	for start < len(lines) && lines[start] == "" {
		start++
	}

	end := len(lines) - 1
	for end >= 0 && lines[end] == "" {
		end--
	}

	if start >= end {
		return nil, fmt.Errorf("banner file contains no content")
	}

	cleanLines := lines[start : end+1]

	currentChar := byte(32)
	index := 0

	for currentChar <= 126 {
		if index+8 > len(cleanLines) {
			return nil, fmt.Errorf("incomplete banner format for character %c at line %d", currentChar, index)
		}

		// Extract 8 lines
		charLines := make([]string, 8)
		for j := 0; j < 8; j++ {
			charLines[j] = cleanLines[index+j]
		}

		asciiMap[currentChar] = charLines
		currentChar++
		index += 9 // Skip the separator line (which might be empty or missing)

		// If we've gone beyond available lines, break
		if index >= len(cleanLines) {
			break
		}
	}

	if len(asciiMap) != 95 {
		return nil, fmt.Errorf("expected 95 characters, got %d", len(asciiMap))
	}

	return asciiMap, nil
}

func readBannerFile(banner string) ([]string, error) {
	data, err := os.ReadFile(banner + ".txt")
	if err != nil {
		return nil, err
	}
	text := string(data)
	lines := strings.Split(text, "\n")

	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	// REMOVED: fmt.Println(lines)
	return lines, nil
}

func main() {
	// Your argument parsing code (keep as is, but remove the \n special case)
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Fprintf(os.Stderr, "Usage: [STRING] [BANNER]\n")
		fmt.Fprintf(os.Stderr, "Banner options: 'standard', 'shadow', 'thinkertoy'\n")
		os.Exit(1)
	}
	// REMOVED the incomplete \n handling

	banner := "standard"
	if len(os.Args) == 3 {
		banner = os.Args[2]
	}

	validBanners := map[string]bool{
		"standard":   true,
		"shadow":     true,
		"thinkertoy": true,
	}
	if !validBanners[banner] {
		fmt.Fprint(os.Stderr, "Invalid banner: use 'standard', 'shadow', or 'thinkertoy'\n")
		os.Exit(1)
	}

	lines, err := readBannerFile(banner)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading banner file: %v\n", err)
		os.Exit(1)
	}

	// Parse banner into map
	asciiMap, err := parseBanner(lines)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing banner: %v\n", err)
		os.Exit(1)
	}

	// Debug output - will remove in Step 3
	fmt.Printf("Successfully loaded %d characters\n", len(asciiMap))

	// Quick test: Print first 3 lines of 'A' if it exists
	if a, ok := asciiMap['A']; ok {
		fmt.Println("First 3 lines of 'A':")
		for i := 0; i < 3 && i < len(a); i++ {
			fmt.Printf("Line %d: %q\n", i, a[i])
		}
	}
}
