package alias

import (
	"os"
	"path/filepath"
)

func FindPathByAlias(target string) (string, bool) {
	aliases, err := ReadAliases()
	if err != nil {
		return "", false
	}
	path, ok := aliases[target]
	return path, ok
}

func Priority(target string) (string, error) {
	if info, err := os.Stat(target); err == nil && info.IsDir() {
		return filepath.Abs(target)
	}

	if path, ok := FindPathByAlias(target); ok {
		return path, nil
	}

	return "", os.ErrNotExist
}
