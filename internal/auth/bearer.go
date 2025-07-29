package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authorization := headers.Get("Authorization")
	if authorization == "" {
		return "", errors.New("authorization header missing")
	}

	if !strings.HasPrefix(authorization, "Bearer ") {
		return "", errors.New("authorization header must start with 'Bearer '")
	}

	token := strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer "))
	if len(token) == 0 {
		return "", errors.New("bearer token is missing")
	}

	return token, nil
}
