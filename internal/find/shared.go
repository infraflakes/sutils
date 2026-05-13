package find

import (
	"fmt"
	"os"

	"github.com/infraflakes/sutils/internal/helper/cli"
	"github.com/infraflakes/sutils/internal/helper/exec"
)

func DeletePath(path string, isDir bool) {
	var err error
	if isDir {
		err = os.RemoveAll(path)
	} else {
		err = os.Remove(path)
	}
	if err != nil {
		fmt.Printf("Failed to delete %s: %v\n", path, err)
	} else {
		fmt.Printf("Deleted: %s\n", path)
	}
}

func FindAndProcess(path string, terms []string, findType string, searchMessage string, deletePrompt string, del bool) {
	isDir := findType == "d"
	for _, term := range terms {
		fmt.Printf(searchMessage, term, path)
		matches, err := exec.RunCommand("find", path, "-type", findType, "-name", fmt.Sprintf("*%s*", term))
		if err != nil {
			fmt.Printf("Error finding paths: %v\n", err)
			continue
		}
		for _, match := range matches {
			fmt.Println(match)
		}

		if del && len(matches) > 0 && cli.Confirm(deletePrompt) {
			for _, match := range matches {
				DeletePath(match, isDir)
			}
		}
	}
}

func DeleteGrepMatches(basePath, term string) {
	matches, err := exec.RunCommand("grep", "-rlE", term, basePath)
	if err != nil {
		fmt.Printf("Error finding files to delete: %v\n", err)
		return
	}
	for _, match := range matches {
		DeletePath(match, false)
	}
}
