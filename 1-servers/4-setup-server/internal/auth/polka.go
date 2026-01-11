package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	s := strings.Split(headers.Get("Authorization"), " ")
	headerValue := strings.TrimSpace(s[0])
	if headerValue != "ApiKey" {
		return "", fmt.Errorf("not contains ApiKey part")
	}
	return strings.TrimSpace(s[1]), nil
}
