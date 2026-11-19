package main

import (
	"encoding/json/jsontext"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStudentsRoute(t *testing.T) {
	t.Parallel()

	//nolint:gosmopolitan // Use the names from the assignment's response example.
	const wantBody = `{
		"students": [{
			"id": 1,
			"name": "佐藤",
			"loginId": "foo123",
			"classroom": {"id": 1, "name": "クラスA"}
		}],
		"totalCount": 1
	}`

	want := jsontext.Value(wantBody)

	err := want.Canonicalize()
	if err != nil {
		t.Fatalf("canonicalize expected JSON: %v", err)
	}

	for _, facilitatorID := range []string{"1", "999"} {
		t.Run(facilitatorID, func(t *testing.T) {
			t.Parallel()

			path := "/students?facilitator_id=" + facilitatorID
			request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, path, nil)
			response := httptest.NewRecorder()
			newRouter().ServeHTTP(response, request)

			if response.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
			}

			const wantContentType = "application/json; charset=utf-8"

			if contentType := response.Header().Get("Content-Type"); contentType != wantContentType {
				t.Errorf("Content-Type = %q, want %q", contentType, wantContentType)
			}

			body := response.Body.String()
			got := jsontext.Value(body)

			err := got.Canonicalize()
			if err != nil {
				t.Fatalf("canonicalize response JSON: %v", err)
			}

			if got.String() != want.String() {
				t.Errorf("response = %s, want %s", body, wantBody)
			}
		})
	}
}

func TestStudentsRouteInvalidFacilitatorID(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		path string
	}{
		{name: "missing", path: "/students"},
		{name: "empty", path: "/students?facilitator_id="},
		{name: "text", path: "/students?facilitator_id=abc"},
		{name: "fraction", path: "/students?facilitator_id=1.5"},
		{name: "zero", path: "/students?facilitator_id=0"},
		{name: "negative", path: "/students?facilitator_id=-1"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, tc.path, nil)
			response := httptest.NewRecorder()
			newRouter().ServeHTTP(response, request)

			if response.Code != http.StatusBadRequest {
				t.Errorf("status = %d, want %d", response.Code, http.StatusBadRequest)
			}

			if response.Body.Len() != 0 {
				t.Errorf("body = %q, want empty", response.Body.String())
			}
		})
	}
}
