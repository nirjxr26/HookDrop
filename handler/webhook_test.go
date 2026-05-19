package handler_test

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/your-username/hookdrop/handler"
	"github.com/your-username/hookdrop/store"
)

func TestHealthz(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/healthz", nil)
	handler.Healthz(rr, req)
	if rr.Code != 200 {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	var body map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("unexpected health body: %#v", body)
	}
}

func TestWebhookPostAndList(t *testing.T) {
	st := store.New()
	h := handler.NewWebhookHandler(st)

	// POST a webhook
	payload := `{"hello":"world"}`
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/h/test", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	h.ServeHTTP(rr, req)
	if rr.Code != 200 {
		t.Fatalf("expected POST status 200, got %d", rr.Code)
	}

	// Verify the store has the event
	events := st.List("test")
	if len(events) != 1 {
		t.Fatalf("expected 1 event in store, got %d", len(events))
	}
	if events[0].Body != payload {
		t.Fatalf("unexpected event body: %q", events[0].Body)
	}
}
