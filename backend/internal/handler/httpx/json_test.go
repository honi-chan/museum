package httpx

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// --------------------------------------------------
// DecodeJSON
// --------------------------------------------------

func TestDecodeJSON(t *testing.T) {
	body := bytes.NewBufferString(
		`{"title":"THINGS I MADE"}`,
	)

	request := httptest.NewRequest(
		http.MethodPost,
		"/exhibitions",
		body,
	)

	var decoded struct {
		Title string `json:"title"`
	}

	err := DecodeJSON(
		request,
		&decoded,
	)

	if err != nil {
		t.Fatalf(
			"DecodeJSON() returned error: %v",
			err,
		)
	}

	if decoded.Title != "THINGS I MADE" {
		t.Errorf(
			"Title = %q, want %q",
			decoded.Title,
			"THINGS I MADE",
		)
	}
}

// --------------------------------------------------
// WriteJSON
// --------------------------------------------------

func TestWriteJSON(t *testing.T) {
	recorder := httptest.NewRecorder()

	response := struct {
		ID string `json:"id"`
	}{
		ID: "exhibition-001",
	}

	err := WriteJSON(
		recorder,
		http.StatusCreated,
		response,
	)

	if err != nil {
		t.Fatalf(
			"WriteJSON() returned error: %v",
			err,
		)
	}

	// Status Codeを確認する。
	if recorder.Code != http.StatusCreated {
		t.Errorf(
			"status = %d, want %d",
			recorder.Code,
			http.StatusCreated,
		)
	}

	// Content-Typeを確認する。
	if recorder.Header().Get(
		"Content-Type",
	) != "application/json" {
		t.Errorf(
			"Content-Type = %q, want %q",
			recorder.Header().Get("Content-Type"),
			"application/json",
		)
	}
}

// --------------------------------------------------
// WriteError
// --------------------------------------------------

func TestWriteError(t *testing.T) {
	recorder := httptest.NewRecorder()

	err := WriteError(
		recorder,
		http.StatusBadRequest,
		errors.New("invalid request"),
	)

	if err != nil {
		t.Fatalf(
			"WriteError() returned error: %v",
			err,
		)
	}

	if recorder.Code != http.StatusBadRequest {
		t.Errorf(
			"status = %d, want %d",
			recorder.Code,
			http.StatusBadRequest,
		)
	}

	expectedBody :=
		"{\"error\":\"invalid request\"}\n"

	if recorder.Body.String() != expectedBody {
		t.Errorf(
			"body = %q, want %q",
			recorder.Body.String(),
			expectedBody,
		)
	}
}
