package auth

import (
	"errors"
	"strings"
	"net/http"
)

func GetAPIKey(headers http.Header) (string, error) {
	authToken := headers.Get("Authorization")
	if authToken == "" {
		return "", errors.New("Unauthorized access")
	}

	elements := strings.Fields(authToken)
	if len(elements) != 2 || elements[0] != "ApiKey" {
		return "", errors.New("Unauthorized access")
	}

	return elements[1], nil
}
