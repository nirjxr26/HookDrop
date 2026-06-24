package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/your-username/hookdrop/store"
)

const (
	// mockStripeIP1 represents a mock client IP address for Stripe webhook seed events.
	mockStripeIP1 = "3.18.12.1"
	// mockStripeIP2 represents a mock client IP address for Stripe webhook seed events.
	mockStripeIP2 = "3.18.12.2"
	// mockGitHubIP1 represents a mock client IP address for GitHub webhook seed events.
	mockGitHubIP1 = "140.82.115.1"
	// mockGitHubIP2 represents a mock client IP address for GitHub webhook seed events.
	mockGitHubIP2 = "140.82.115.2"
	// mockAuth0IP represents a mock client IP address for Auth0 webhook seed events.
	mockAuth0IP   = "54.221.18.44"
)

func seedMockData(s *store.Store) {
	now := time.Now().UTC()

	// ---- STRIPE BUCKET EVENTS ----
	s.Add("stripe", store.WebhookEvent{
		TraceID:   uuid.NewString(),
		BucketID:  "stripe",
		Method:    "POST",
		SourceIP:  mockStripeIP1,
		Timestamp: now.Add(-2 * time.Minute),
		Headers: map[string][]string{
			"Content-Type":     {"application/json"},
			"User-Agent":       {"Stripe/1.0 (+https://stripe.com/docs/webhooks)"},
			"Stripe-Signature": {"t=1718536800,v1=9fa8db31c26ebef8781a7b8e1a7b8e1a7b8e1a7b8e1a7b8e1a7b8e1a7b8e1a7b"},
		},
		Body: `{
  "id": "evt_1Oc9c2HzG65Vn8uS7d8f9g0h",
  "object": "event",
  "api_version": "2023-10-16",
  "created": 1718536800,
  "type": "charge.succeeded",
  "data": {
    "object": {
      "id": "ch_3Oc9c2HzG65Vn8uS1a2b3c4d",
      "object": "charge",
      "amount": 2999,
      "amount_captured": 2999,
      "amount_refunded": 0,
      "billing_details": {
        "address": {
          "city": "San Francisco",
          "country": "US",
          "line1": "123 Market St",
          "postal_code": "94105",
          "state": "CA"
        },
        "email": "customer@example.com",
        "name": "Jane Doe"
      },
      "captured": true,
      "currency": "usd",
      "customer": "cus_Oz8f9g0h1a2b",
      "description": "Subscription to Premium Plan",
      "paid": true,
      "status": "succeeded"
    }
  }
}`,
	})

	s.Add("stripe", store.WebhookEvent{
		TraceID:   uuid.NewString(),
		BucketID:  "stripe",
		Method:    "POST",
		SourceIP:  mockStripeIP2,
		Timestamp: now.Add(-15 * time.Minute),
		Headers: map[string][]string{
			"Content-Type":     {"application/json"},
			"User-Agent":       {"Stripe/1.0 (+https://stripe.com/docs/webhooks)"},
			"Stripe-Signature": {"t=1718536000,v1=a1b2c3d4e5f6g7h8i9j0a1b2c3d4e5f6g7h8i9j0a1b2c3d4e5f6g7h8i9j0a1b2"},
		},
		Body: `{
  "id": "evt_1Oc9b2HzG65Vn8uS5a6b7c8d",
  "object": "event",
  "api_version": "2023-10-16",
  "created": 1718536000,
  "type": "customer.subscription.created",
  "data": {
    "object": {
      "id": "sub_1Oc9b2HzG65Vn8uS9d0e1f2g",
      "object": "subscription",
      "customer": "cus_Oz8f9g0h1a2b",
      "status": "active",
      "current_period_start": 1718536000,
      "current_period_end": 1721128000,
      "items": {
        "object": "list",
        "data": [
          {
            "id": "si_Oz8f9g0h2c3d",
            "object": "subscription_item",
            "price": {
              "id": "price_1Oc9b2HzG65Vn8uS123abcde",
              "product": "prod_PremiumSubscription",
              "unit_amount": 2999,
              "currency": "usd"
            }
          }
        ]
      }
    }
  }
}`,
	})

	// ---- GITHUB BUCKET EVENTS ----
	s.Add("github", store.WebhookEvent{
		TraceID:   uuid.NewString(),
		BucketID:  "github",
		Method:    "POST",
		SourceIP:  mockGitHubIP1,
		Timestamp: now.Add(-5 * time.Minute),
		Headers: map[string][]string{
			"Content-Type":          {"application/json"},
			"User-Agent":            {"GitHub-Hookshot/7a82fb"},
			"X-GitHub-Event":        {"push"},
			"X-GitHub-Delivery":     {"d8a2f64c-bcae-11ee-8e8e-c34d399634d3"},
			"X-Hub-Signature-256":   {"sha256=1bc2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2"},
		},
		Body: `{
  "ref": "refs/heads/main",
  "before": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0",
  "after": "9z8y7x6w5v4u3t2s1r0q9p8o7n6m5l4k3j2i1h0g",
  "repository": {
    "id": 84729184,
    "name": "Polaris",
    "full_name": "nirjxr26/Polaris",
    "private": true,
    "owner": {
      "name": "nirjxr26",
      "email": "nirjxr26@users.noreply.github.com"
    },
    "html_url": "https://github.com/nirjxr26/Polaris"
  },
  "pusher": {
    "name": "nirjxr26",
    "email": "nirjxr26@users.noreply.github.com"
  },
  "commits": [
    {
      "id": "9z8y7x6w5v4u3t2s1r0q9p8o7n6m5l4k3j2i1h0g",
      "message": "feat: upgrade real-time dashboard UI",
      "timestamp": "2026-06-16T04:22:00Z",
      "author": {
        "name": "Antigravity",
        "email": "antigravity@gemini.ai"
      }
    }
  ]
}`,
	})

	s.Add("github", store.WebhookEvent{
		TraceID:   uuid.NewString(),
		BucketID:  "github",
		Method:    "POST",
		SourceIP:  mockGitHubIP2,
		Timestamp: now.Add(-30 * time.Minute),
		Headers: map[string][]string{
			"Content-Type":          {"application/json"},
			"User-Agent":            {"GitHub-Hookshot/7a82fb"},
			"X-GitHub-Event":        {"pull_request"},
			"X-GitHub-Delivery":     {"e9b3f75d-cdaf-11ee-8e8e-d45e400745e4"},
			"X-Hub-Signature-256":   {"sha256=2cd3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3"},
		},
		Body: `{
  "action": "opened",
  "number": 42,
  "pull_request": {
    "id": 1492049102,
    "title": "docs: update API documentation references",
    "user": {
      "login": "developer-adam",
      "id": 551920
    },
    "body": "This PR resolves the missing endpoints list in the README.",
    "created_at": "2026-06-16T03:55:00Z",
    "head": {
      "ref": "docs/readme-endpoints",
      "sha": "f1e2d3c4b5a697887766554433221100"
    },
    "base": {
      "ref": "main"
    }
  }
}`,
	})

	// ---- AUTH0 BUCKET EVENTS ----
	s.Add("auth0", store.WebhookEvent{
		TraceID:   uuid.NewString(),
		BucketID:  "auth0",
		Method:    "POST",
		SourceIP:  mockAuth0IP,
		Timestamp: now.Add(-8 * time.Minute),
		Headers: map[string][]string{
			"Content-Type":     {"application/json"},
			"User-Agent":       {"Auth0-Hook-Worker"},
			"X-Auth0-Event-Id": {"evt_auth0_987452"},
		},
		Body: `{
  "action": "user-signup",
  "tenant": "dev-auth-polaris",
  "user": {
    "user_id": "auth0|666ecd9f3a218d007123abcd",
    "email": "newuser@example.com",
    "email_verified": false,
    "username": "newuser2026",
    "created_at": "2026-06-16T04:24:00.123Z",
    "updated_at": "2026-06-16T04:24:00.123Z",
    "user_metadata": {
      "first_name": "John",
      "last_name": "Smith"
    }
  },
  "strategy": "auth0",
  "connection": "Username-Password-Authentication"
}`,
	})

	// ---- DEFAULT BUCKET EVENTS ----
	s.Add("default", store.WebhookEvent{
		TraceID:   uuid.NewString(),
		BucketID:  "default",
		Method:    "POST",
		SourceIP:  "127.0.0.1",
		Timestamp: now.Add(-10 * time.Second),
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
			"User-Agent":   {"curl/8.4.0"},
		},
		Body: `{
  "message": "Welcome to HookDrop! Send any payload to this endpoint to test."
}`,
	})
}
