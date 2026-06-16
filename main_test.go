package main

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	err := os.Setenv("TEST_ENV_VAR", "hello")
	if err != nil {
		t.Fatalf("failed to set env: %v", err)
	}
	defer func() {
		_ = os.Unsetenv("TEST_ENV_VAR")
	}()

	val := getEnv("TEST_ENV_VAR", "default")
	if val != "hello" {
		t.Errorf("expected hello, got %s", val)
	}

	valDefault := getEnv("NON_EXISTENT_VAR", "fallback")
	if valDefault != "fallback" {
		t.Errorf("expected fallback, got %s", valDefault)
	}
}
