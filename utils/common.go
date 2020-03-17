package utils

import (
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
