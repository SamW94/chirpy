package auth

import (
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	return strings.Replace(headers.Get("Authorization"), "ApiKey ", "", -1), nil
}
