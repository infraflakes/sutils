package alias

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/infraflakes/sutils/internal/helper/fs"
)

const ConfigFileName = "scd-alias.conf"

func getConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(configDir, "scd", "scd-alias.conf")
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", err
	}
	return path, nil
}

func ReadAliases() (map[string]string, error) {
	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	aliases := make(map[string]string)
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return aliases, nil
		}
		return nil, err
	}
	defer fs.CloseFile(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			aliases[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return aliases, scanner.Err()
}

func SaveAliases(aliases map[string]string) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	file, err := fs.CreateFile(path)
	if err != nil {
		return err
	}
	defer fs.CloseFile(file)

	writer := bufio.NewWriter(file)
	for alias, p := range aliases {
		_, err := fmt.Fprintf(writer, "%s = %s\n", alias, p)
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

func AddAlias(alias string, path string) error {
	aliases, err := ReadAliases()
	if err != nil {
		return err
	}
	aliases[alias] = path
	return SaveAliases(aliases)
}

func RemoveAlias(alias string) error {
	aliases, err := ReadAliases()
	if err != nil {
		return err
	}
	if _, ok := aliases[alias]; !ok {
		return fmt.Errorf("alias '%s' not found", alias)
	}
	delete(aliases, alias)
	return SaveAliases(aliases)
}

func WipeAliases() error {
	return SaveAliases(make(map[string]string))
}

func ExportAliases(destPath string) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("no aliases available to export")
		}
		return err
	}

	return os.WriteFile(destPath, data, 0755)
}
