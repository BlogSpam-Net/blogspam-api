//
// Simple testing of the HTTP-server
//
//
package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

//
// Submitting JSON must be done via a POST.
//
func TestTestMethod(t *testing.T) {

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SpamTestHandler)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Unexpected status-code: %v", status)
	}

	// Check the response body is what we expect.
	expected := "Must be called via HTTP-POST\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got '%v' want '%v'",
			rr.Body.String(), expected)
	}
}

//
// Test that we can POST to the end-point
//
func TestSpam(t *testing.T) {
	body := []byte("{\"comment\":\"Moi Kissa\",\"name\":\"http://example.com\"}")

	req, err := http.NewRequest("POST", "/", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SpamTestHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Unexpected status-code: %v", status)
	}

	if !strings.Contains(rr.Body.String(), "SPAM") {
		t.Errorf("Body was '%v' not spam",
			rr.Body.String())
	}
}

//
// Test that we can handle bogus JSON POSTed to the end-point
//
func TestBogusJSON(t *testing.T) {
	body := []byte("{\"comment\",\"Moi Kissa\",\"name\":\"http://example.com\"}")

	req, err := http.NewRequest("POST", "/", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SpamTestHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Unexpected status-code: %v", status)
	}

	if !strings.Contains(rr.Body.String(), "invalid character") {
		t.Errorf("Body was '%v' not a bogus JSON error",
			rr.Body.String())
	}
}

//
// The name-check plugin will mark a field as spam, so we'll test that
// we can exclude that.
//
func TestSpamExclusion(t *testing.T) {
	body := []byte("{\"options\":\"exclude=name\",\"comment\":\"Moi Kissa\",\"name\":\"http://example.com\"}")

	req, err := http.NewRequest("POST", "/", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SpamTestHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Unexpected status-code: %v", status)
	}

	if !strings.Contains(rr.Body.String(), "OK") {
		t.Errorf("Body was '%v' not OK",
			rr.Body.String())
	}
}

//
// Test that we can retrieve content from the plugins-list end-point
//
func TestPluginList(t *testing.T) {

	//
	// Make the request
	//
	req, err := http.NewRequest("GET", "/list/", nil)
	if err != nil {
		t.Fatal(err)
	}

	//
	// Fake it out
	//
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(PluginListHandler)
	handler.ServeHTTP(rr, req)

	//
	// Test the status-code is OK
	//
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Unexpected status-code: %v", status)
	}

	//
	// Test that the body contained our expected content.
	//
	if !strings.Contains(rr.Body.String(), "requiremx") {
		t.Fatalf("Unexpected body: '%s'", rr.Body.String())
	}

}
