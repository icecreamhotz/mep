package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func SplitTokenFromHeader(token string) (string, bool) {
	splitToken := strings.Split(token, "Bearer")
	if len(splitToken) != 2 {
		return "", false
	}
	token = strings.TrimSpace(splitToken[1])

	return token, true
}

func GetBaseDirectory() (string, error) {
	baseDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	return baseDir, err
}
