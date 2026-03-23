package config

import (
	"errors"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func Load(path string) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil
	}

	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	return godotenv.Load(path)
}
