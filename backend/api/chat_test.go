package handler

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestHandler_RejectsNonPostMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/chat", nil)
	rr := httptest.NewRecorder()

	Handler(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
}

func TestHandler_RejectsInvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/chat", strings.NewReader("{"))
	rr := httptest.NewRecorder()

	Handler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestHandler_RejectsEmptyMessage(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/chat", strings.NewReader(`{"message":"   "}`))
	rr := httptest.NewRecorder()

	Handler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestHandler_RequiresAPIKey(t *testing.T) {
	_ = os.Unsetenv("DEEPSEEK_API_KEY")
	req := httptest.NewRequest(http.MethodPost, "/api/chat", strings.NewReader(`{"message":"hello"}`))
	rr := httptest.NewRecorder()

	Handler(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}

func TestHandler_OptionsPreflight(t *testing.T) {
	req := httptest.NewRequest(http.MethodOptions, "/api/chat", nil)
	rr := httptest.NewRecorder()

	Handler(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rr.Code)
	}
}
