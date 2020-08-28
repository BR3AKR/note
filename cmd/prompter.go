package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var s *bufio.Scanner

func prompt(prompt string) string {
	fmt.Print(prompt)
	if s == nil {
		s = bufio.NewScanner(os.Stdin)
	}
	s.Scan()
	return strings.TrimSpace(s.Text())
}

func repeatPrompt(promptStr string) []string {
	values := make([]string, 0)

	for c := prompt(promptStr); c != ""; c = prompt(promptStr) {
		values = append(values, c)
	}

	return values
}
