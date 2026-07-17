package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIKeyAuth(t *testing.T) {
	tests := []struct {
		name       string
		header     string
		wantStatus int
		wantNext   bool
	}{
		{"missing api key", "", http.StatusUnauthorized, false},
		{"invalid api key", "wrong", http.StatusUnauthorized, false},
		{"valid api key", "secret", http.StatusOK, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextCalled := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusOK)
			})
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.header != "" {
				req.Header.Set("X-API-Key", tt.header)
			}
			rec := httptest.NewRecorder()

			APIKeyAuth("secret")(next).ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if nextCalled != tt.wantNext {
				t.Fatalf("nextCalled = %t, want %t", nextCalled, tt.wantNext)
			}
		})
	}
}
