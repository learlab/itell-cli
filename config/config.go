package config

import (
	"os"
	"path"
)

func GetDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		return path.Join(os.TempDir(), ".itell")
	}
	return path.Join(dir, ".itell")
}
