package httpapi

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	authusecase "codelife-study-be/internal/usecase/auth"
)

func TestDocumentDetailRequiresAuthentication(t *testing.T) {
	authService := authusecase.New(nil, nil, "test-secret", time.Minute, time.Hour)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := New(nil, authService, nil, nil, nil, logger, 1024)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/documents/private-document", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", response.Code)
	}
}
