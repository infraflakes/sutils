package fs

import (
	"fmt"
	"os"
)

func CreateFile(path string) (*os.File, error) {
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return nil, err
	}
	return file, nil
}

func OpenFile(path string) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	return file, nil
}

func CloseFile(file *os.File) {
	if err := file.Close(); err != nil {
		fmt.Println("Error closing file:", err)
	}
}

func LogError(file *os.File, message string) {
	if _, err := fmt.Fprintln(file, message); err != nil {
		fmt.Println("Error writing to log file:", err)
	}
}
