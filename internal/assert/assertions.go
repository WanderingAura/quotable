package assert

import (
	"strings"
	"testing"
)

func Equal[T comparable](t *testing.T, value, expected T) {
	t.Helper()
	if value != expected {
		t.Errorf("expected: %v; got: %v", expected, value)
	}
}

func StringContains(t *testing.T, str, contains string) {
	t.Helper()
	if !strings.Contains(str, contains) {
		t.Errorf("expected \"%s\" to contain \"%s\"", str, contains)
	}
}
