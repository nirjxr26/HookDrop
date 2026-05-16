package store

import (
	"sync"
	"time"
)

type WebhookEvent struct {
	TraceID   string              `json:"trace_id"`
	BucketID  string              `json:"bucket_id,omitempty"`
	Method    string              `json:"method"`
	Headers   map[string][]string `json:"headers"`
	Body      string              `json:"body"`
	SourceIP  string              `json:"source_ip"`
	Timestamp time.Time           `json:"timestamp"`
}

type Bucket struct {
	Events []WebhookEvent
}

type Store struct {
	mu          sync.RWMutex
	buckets     map[string]*Bucket
	subscribers map[string][]chan WebhookEvent
}

const maxBucketEvents = 50

func New() *Store {
	return &Store{
		buckets:     make(map[string]*Bucket),
		subscribers: make(map[string][]chan WebhookEvent),
	}
}

func (s *Store) Add(bucketID string, event WebhookEvent) {
	s.mu.Lock()
	bucket := s.ensureBucketLocked(bucketID)
	if len(bucket.Events) >= maxBucketEvents {
		copy(bucket.Events, bucket.Events[1:])
		bucket.Events = bucket.Events[:maxBucketEvents-1]
	}
	bucket.Events = append(bucket.Events, event)
	subscribers := append([]chan WebhookEvent(nil), s.subscribers[bucketID]...)
	s.mu.Unlock()

	for _, ch := range subscribers {
		select {
		case ch <- event:
		default:
		}
	}
}

func (s *Store) List(bucketID string) []WebhookEvent {
	s.mu.RLock()
	bucket := s.buckets[bucketID]
	if bucket == nil {
		s.mu.RUnlock()
		return []WebhookEvent{}
	}
	events := make([]WebhookEvent, len(bucket.Events))
	copy(events, bucket.Events)
	s.mu.RUnlock()
	return events
}

func (s *Store) Subscribe(bucketID string) chan WebhookEvent {
	ch := make(chan WebhookEvent, 16)
	s.mu.Lock()
	s.ensureBucketLocked(bucketID)
	s.subscribers[bucketID] = append(s.subscribers[bucketID], ch)
	s.mu.Unlock()
	return ch
}

func (s *Store) Unsubscribe(bucketID string, ch chan WebhookEvent) {
	s.mu.Lock()
	subscribers := s.subscribers[bucketID]
	if len(subscribers) == 0 {
		s.mu.Unlock()
		return
	}
	for i, subscriber := range subscribers {
		if subscriber == ch {
			s.subscribers[bucketID] = append(subscribers[:i], subscribers[i+1:]...)
			break
		}
	}
	s.mu.Unlock()
}

func (s *Store) ensureBucketLocked(bucketID string) *Bucket {
	bucket := s.buckets[bucketID]
	if bucket == nil {
		bucket = &Bucket{}
		s.buckets[bucketID] = bucket
	}
	return bucket
}
