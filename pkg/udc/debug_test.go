package udc

import (
	"testing"
)

func TestDebugMode(t *testing.T) {
	// Test initial state
	if IsDebugMode() {
		t.Error("Expected debug mode to be false initially")
	}

	// Test enabling debug mode
	SetDebugMode(true)
	if !IsDebugMode() {
		t.Error("Expected debug mode to be true after enabling")
	}

	// Test disabling debug mode
	SetDebugMode(false)
	if IsDebugMode() {
		t.Error("Expected debug mode to be false after disabling")
	}

	// Test enabling again
	SetDebugMode(true)
	if !IsDebugMode() {
		t.Error("Expected debug mode to be true after re-enabling")
	}

	// Clean up - disable debug mode
	SetDebugMode(false)
}
