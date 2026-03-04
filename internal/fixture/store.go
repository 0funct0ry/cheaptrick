package fixture

import (
	"os"
	"path/filepath"
)

func GetFixture(dir, hash string) (string, bool) {
	if dir == "" {
		return "", false
	}
	path := filepath.Join(dir, hash+".json")
	b, err := os.ReadFile(path)
	if err == nil {
		return string(b), true
	}
	return "", false
}

func SaveFixture(dir, hash, response string) error {
	if dir == "" {
		return nil
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	path := filepath.Join(dir, hash+".json")
	return os.WriteFile(path, []byte(response), 0644)
}
