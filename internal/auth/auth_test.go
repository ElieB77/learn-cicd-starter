package auth

import (
	"errors"
	"net/http"
	"testing"
)

func TestGetAPIKey(t *testing.T) {
	// Test case: Valid API Key
	headers := http.Header{
		"Authorization": []string{"ApiKey my-secret-key"},
	}
	got, err := GetAPIKey(headers)
	want := "my-secret-key"
	if err != nil || got != want {
		t.Fatalf("expected: %v, got: %v, error: %v", want, got, err)
	}

	// Test case: Missing Authorization Header
	headers = http.Header{}
	_, err = GetAPIKey(headers)
	if err == nil || err != ErrNoAuthHeaderIncluded {
		t.Fatalf("expected error: %v, got: %v", ErrNoAuthHeaderIncluded, err)
	}

	// Test case: Malformed Authorization Header - No ApiKey prefix
	headers = http.Header{
		"Authorization": []string{"Bearer my-secret-key"},
	}
	_, err = GetAPIKey(headers)
	wantErr := errors.New("malformed authorization header")
	if err == nil || err.Error() != wantErr.Error() {
		t.Fatalf("expected error: %v, got: %v", wantErr, err)
	}

	// Test case: Malformed Authorization Header - Missing key
	headers = http.Header{
		"Authorization": []string{"ApiKey"},
	}
	_, err = GetAPIKey(headers)
	if err == nil || err.Error() != wantErr.Error() {
		t.Fatalf("expected error: %v, got: %v", wantErr, err)
	}
}
