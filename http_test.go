package main

import (
	"encoding/json/jsontext"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetStudentsHTTP(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		path       string
		wantStatus int
		assertBody func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "valid", path: "/students?facilitator_id=1",
			wantStatus: http.StatusOK, assertBody: assertStudentsJSON,
		},
		{
			name: "unknown", path: "/students?facilitator_id=999",
			wantStatus: http.StatusOK, assertBody: assertStudentsJSON,
		},
		{
			name: "missing", path: "/students",
			wantStatus: http.StatusBadRequest, assertBody: assertEmptyBody,
		},
		{
			name: "empty", path: "/students?facilitator_id=",
			wantStatus: http.StatusBadRequest, assertBody: assertEmptyBody,
		},
		{
			name: "text", path: "/students?facilitator_id=abc",
			wantStatus: http.StatusBadRequest, assertBody: assertEmptyBody,
		},
		{
			name: "fraction", path: "/students?facilitator_id=1.5",
			wantStatus: http.StatusBadRequest, assertBody: assertEmptyBody,
		},
		{
			name: "zero", path: "/students?facilitator_id=0",
			wantStatus: http.StatusBadRequest, assertBody: assertEmptyBody,
		},
		{
			name: "negative", path: "/students?facilitator_id=-1",
			wantStatus: http.StatusBadRequest, assertBody: assertEmptyBody,
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			router := newRouter()
			request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, testCase.path, nil)
			response := httptest.NewRecorder()

			// Act
			router.ServeHTTP(response, request)

			// Assert
			if response.Code != testCase.wantStatus {
				t.Fatalf("status = %d, want %d", response.Code, testCase.wantStatus)
			}

			testCase.assertBody(t, response)
		})
	}
}

func assertStudentsJSON(t *testing.T, response *httptest.ResponseRecorder) {
	t.Helper()

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

	const wantContentType = "application/json; charset=utf-8"

	if contentType := response.Header().Get("Content-Type"); contentType != wantContentType {
		t.Errorf("Content-Type = %q, want %q", contentType, wantContentType)
	}

	body := response.Body.String()
	got := jsontext.Value(body)

	err = got.Canonicalize()
	if err != nil {
		t.Fatalf("canonicalize response JSON: %v", err)
	}

	if got.String() != want.String() {
		t.Errorf("response = %s, want %s", body, wantBody)
	}
}

func assertEmptyBody(t *testing.T, response *httptest.ResponseRecorder) {
	t.Helper()

	if response.Body.Len() != 0 {
		t.Errorf("body = %q, want empty", response.Body.String())
	}
}
