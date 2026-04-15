package testutil

import "testing"

func AssertNilError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("expected error to be nil, but got %v", err)
	}
}

func AssertStatusCode(t *testing.T, got, want int) {
	t.Helper()

	if got != want {
		t.Fatalf("expected response.StatusCode to be [%d], but got [%d]", want, got)
	}
}