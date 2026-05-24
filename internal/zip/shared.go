package zip

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/infraflakes/sutils/internal/helper/exec"
	"github.com/infraflakes/sutils/internal/helper/utils"
)

func ExpandTargets(targets []string) []string {
	for i, target := range targets {
		info, err := os.Stat(target)
		if err == nil && info.IsDir() {
			targets[i] = filepath.Join(target, "/*")
		}
	}
	return targets
}

func BuildArchiveCommand(archiveName string, targets []string, password string) {
	fileExt := filepath.Ext(archiveName)
	cmdArgs := []string{"a"}

	if password != "" {
		cmdArgs = append(cmdArgs, "-p"+password)
	}

	if fileExt == ".zip" {
		cmdArgs = append(cmdArgs, "-tzip")
	}

	cmdArgs = append(cmdArgs, archiveName)
	cmdArgs = append(cmdArgs, targets...)

	utils.CheckErr(run7z(cmdArgs, password))
}

func BuildExtractCommand(target string, password string) {
	cmdArgs := []string{"x", "-aoa"}
	if password != "" {
		cmdArgs = append(cmdArgs, "-p"+password)
	}
	cmdArgs = append(cmdArgs, target)
	utils.CheckErr(run7z(cmdArgs, password))
}

func run7z(args []string, password string) error {
	err := exec.ExecuteCommand("7z", args...)
	if err != nil && password != "" {
		return fmt.Errorf("%s", strings.ReplaceAll(err.Error(), "-p"+password, "-p***"))
	}
	return err
}
