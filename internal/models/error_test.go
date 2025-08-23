package models

import "testing"


func TestNewErrorResponse(t *testing.T) {
	err := NewErrorResponse("Test Error", 400, "Test detail", "/test")

	if err.Type != "about:blank" {
		t.Errorf("Expected Type to be 'about:blank', got %s", err.Type)
	}

	if err.Title != "Test Error" {
		t.Errorf("Expected Title to be 'Test Error', got %s", err.Title)
	}

	if err.Status != 400 {
		t.Errorf("Expected Status to be 400, got %d", err.Status)
	}

	if err.Detail != "Test detail" {
		t.Errorf("Expected Detail to be 'Test detail', got %s", err.Detail)
	}

	if err.Instance != "/test" {
		t.Errorf("Expected Instance to be '/test', got %s", err.Instance)
	}
}

func TestErrBadRequest(t *testing.T) {
	err := ErrBadRequest("Invalid input", "/test/bad-request")

	if err.Status != 400 {
		t.Errorf("Expected Status to be 400, got %d", err.Status)
	}

	if err.Title != "Bad Request" {
		t.Errorf("Expected Title to be 'Bad Request', got %s", err.Title)
	}
}

func TestErrNotFound(t *testing.T) {
	err := ErrNotFound("User", 123, "/test/not-found")

	if err.Status != 404 {
		t.Errorf("Expected Status to be 404, got %d", err.Status)
	}

	if err.Title != "User Not Found" {
		t.Errorf("Expected Title to be 'User Not Found', got %s", err.Title)
	}
}

func TestErrInternalServer(t *testing.T) {
	err := ErrInternalServer("Database error", "/test/internal-error")

	if err.Status != 500 {
		t.Errorf("Expected Status to be 500, got %d", err.Status)
	}

	if err.Title != "Internal Server Error" {
		t.Errorf("Expected Title to be 'Internal Server Error', got %s", err.Title)
	}
}
