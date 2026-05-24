package alias

import (
	"os"
	"path/filepath"
)

func FindPathByAlias(target string) (string, bool, error) {
	aliases, err := ReadAliases()
	if err != nil {
		return "", false, err
	}
	path, ok := aliases[target]
	return path, ok, nil
}
func Priority(target string) (string, error) {
	if info, err := os.Stat(target); err == nil && info.IsDir() {
		return filepath.Abs(target)
	}
	if path, ok, err := FindPathByAlias(target); err != nil {
		return "", err
	} else if ok {
		return path, nil
	}
	return "", os.ErrNotExist
}
