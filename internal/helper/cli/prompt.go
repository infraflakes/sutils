package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func Confirm(prompt string) bool {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		if errors.Is(err, io.EOF) && strings.TrimSpace(response) != "" {
			return strings.ToLower(strings.TrimSpace(response)) == "y"
		}
		fmt.Printf("Failed to read input: %v\n", err)
		return false
	}
	return strings.ToLower(strings.TrimSpace(response)) == "y"
}

func GetInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		if errors.Is(err, io.EOF) && strings.TrimSpace(input) != "" {
			return strings.TrimSpace(input)
		}
		fmt.Printf("Failed to read input: %v\n", err)
		return ""
	}
	return strings.TrimSpace(input)
}
