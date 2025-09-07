package health

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	Handler(w, req)

	expected := "Erebus is alive 🚀" // no newline
	if w.Body.String() != expected {
		t.Errorf("expected %q, got %q", expected, w.Body.String())
	}
}
