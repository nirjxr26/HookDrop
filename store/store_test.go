package store_test

import (
	"testing"
	"time"

	"github.com/your-username/hookdrop/store"
)

func TestStoreAddAndList(t *testing.T) {
	s := store.New()
	event := store.WebhookEvent{
		TraceID:   "123",
		BucketID:  "test",
		Method:    "POST",
		Body:      `{"data": "val"}`,
		SourceIP:  "127.0.0.1",
		Timestamp: time.Now(),
	}

	s.Add("test", event)

	events := s.List("test")
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].TraceID != "123" {
		t.Errorf("expected trace id 123, got %s", events[0].TraceID)
	}

	// Test max events capping (maxBucketEvents is 50)
	for i := 0; i < 60; i++ {
		s.Add("test", event)
	}
	events = s.List("test")
	if len(events) != 50 {
		t.Errorf("expected capped size of 50 events, got %d", len(events))
	}
}

func TestStoreSubscribe(t *testing.T) {
	s := store.New()
	ch := s.Subscribe("test")
	defer s.Unsubscribe("test", ch)

	event := store.WebhookEvent{
		TraceID: "456",
	}

	s.Add("test", event)

	select {
	case received := <-ch:
		if received.TraceID != "456" {
			t.Errorf("expected trace id 456, got %s", received.TraceID)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for event")
	}
}
