package auth

import (
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	return strings.Replace(headers.Get("Authorization"), "Bearer ", "", -1), nil
}
