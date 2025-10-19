package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// mock server of game services
func mockService(status int, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Write([]byte(body))
	}))
}

func TestCardsRouteForwarding(t *testing.T) {
	mtg_server := mockService(200, `MTG OK`)
	lorcana_server := mockService(200, `Lorcana OK`)

	defer mtg_server.Close()
	defer lorcana_server.Close()

	router_config := Config{
		LorcanaServiceURL: lorcana_server.URL,
		MTGServiceURL:     mtg_server.URL,
	}

	r := NewRouter(router_config)

	tests := []struct {
		game         string
		expectedBody string
	}{
		{"mtg", "MTG OK"},
		{"lorcana", "Lorcana OK"},
	}

	for _, tt := range tests {
		req, _ := http.NewRequest("GET", "/cards?game="+tt.game, nil)
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rr.Code)
		}
		if rr.Body.String() != tt.expectedBody {
			t.Errorf("expected body %q, got %q", tt.expectedBody, rr.Body.String())
		}
	}
}
