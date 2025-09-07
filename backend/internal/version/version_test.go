package version

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	w := httptest.NewRecorder()

	Handler(w, req) // must match actual function name in version.go

	expected := "v0.0.1\n" // keep newline if handler uses Fprintln
	if w.Body.String() != expected {
		t.Errorf("expected %q, got %q", expected, w.Body.String())
	}
}
