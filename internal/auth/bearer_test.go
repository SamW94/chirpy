package auth

import (
	"testing"

	"net/http"
)

func TestGetBearerToken_Success(t *testing.T) {
	httpHeader := http.Header{}
	http.Header.Add(httpHeader, "Authorization", "Bearer cKSXez/B/cjs8h/VcH76NRAKx+7MxudmThYMhv3gEEIwExTcP5h9WfUUBiyRvFQfFnCwtabLRVdNV2KENRgJxw==")

	token, err := GetBearerToken(httpHeader)
	if err != nil {
		t.Fatalf("GetBearerToken failed: %v", err)
	}
	if token != "cKSXez/B/cjs8h/VcH76NRAKx+7MxudmThYMhv3gEEIwExTcP5h9WfUUBiyRvFQfFnCwtabLRVdNV2KENRgJxw==" {
		t.Errorf("Expected cKSXez/B/cjs8h/VcH76NRAKx+7MxudmThYMhv3gEEIwExTcP5h9WfUUBiyRvFQfFnCwtabLRVdNV2KENRgJxw==, got %s", token)
	}
}
