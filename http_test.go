package main

import (
	"database/sql"
	_ "embed"
	"encoding/json/jsontext"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

//go:embed testdata/students.sql
var studentFixtures string

//nolint:paralleltest // Seeded cases insert the same rows and must run sequentially.
func TestGetStudentsHTTP(t *testing.T) {
	pool := openHTTPTestPool(t)

	t.Run("assigned students", func(t *testing.T) {
		const want = `{
			"students": [
				{"id": 1, "name": "佐藤", "loginId": "foo123", "classroom": {"id": 1, "name": "クラスA"}},
				{"id": 4, "name": "加藤", "loginId": "hoge0000", "classroom": {"id": 1, "name": "クラスA"}},
				{"id": 7, "name": "小田原", "loginId": "fizz9999", "classroom": {"id": 1, "name": "クラスA"}},
				{"id": 10, "name": "Williams", "loginId": "xxx:42", "classroom": {"id": 1, "name": "クラスA"}}
			],
			"totalCount": 4
		}`

		response := getStudents(t, newSeededHTTPTestUseCase(t, pool), "/students?facilitator_id=1", http.StatusOK)
		assertStudentsJSON(t, response, want)
	})

	t.Run("unknown facilitator", func(t *testing.T) {
		response := getStudents(t, newSeededHTTPTestUseCase(t, pool), "/students?facilitator_id=999", http.StatusNotFound)
		assertEmptyBody(t, response)
	})

	t.Run("internal server error", func(t *testing.T) {
		transaction, err := pool.BeginTx(t.Context(), nil)
		if err != nil {
			t.Fatalf("begin test transaction: %v", err)
		}

		// A closed transaction makes the repository query fail.
		err = transaction.Rollback()
		if err != nil {
			t.Fatalf("close test transaction: %v", err)
		}

		useCase := newStudentUseCase(newStudentRepository(transaction))
		response := getStudents(t, useCase, "/students?facilitator_id=1", http.StatusInternalServerError)
		assertEmptyBody(t, response)
	})
}

func TestGetStudentsHTTPInvalidFacilitatorID(t *testing.T) {
	t.Parallel()

	pool := openHTTPTestPool(t)

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
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			useCase := newStudentUseCase(newStudentRepository(pool))
			response := getStudents(t, useCase, testCase.path, http.StatusBadRequest)
			assertEmptyBody(t, response)
		})
	}
}

func openHTTPTestPool(t *testing.T) *sql.DB {
	t.Helper()

	config := mysqlConfigFromEnv()
	config.DBName = "compass_quiz_test"
	config.MultiStatements = true

	pool, err := openMySQL(config)
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}

	t.Cleanup(func() {
		err := pool.Close()
		if err != nil {
			t.Errorf("close test database: %v", err)
		}
	})

	return pool
}

func newSeededHTTPTestUseCase(t *testing.T, pool *sql.DB) studentUseCase {
	t.Helper()

	transaction, err := pool.BeginTx(t.Context(), nil)
	if err != nil {
		t.Fatalf("begin test transaction: %v", err)
	}

	t.Cleanup(func() {
		err := transaction.Rollback()
		if err != nil && !errors.Is(err, sql.ErrTxDone) {
			t.Errorf("roll back test transaction: %v", err)
		}
	})

	_, err = transaction.ExecContext(t.Context(), studentFixtures)
	if err != nil {
		t.Fatalf("seed test database: %v", err)
	}

	return newStudentUseCase(newStudentRepository(transaction))
}

func getStudents(t *testing.T, useCase studentUseCase, path string, wantStatus int) *httptest.ResponseRecorder {
	t.Helper()

	request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, path, nil)
	response := httptest.NewRecorder()
	newRouter(useCase).ServeHTTP(response, request)

	if response.Code != wantStatus {
		t.Fatalf("status = %d, want %d", response.Code, wantStatus)
	}

	return response
}

func assertStudentsJSON(t *testing.T, response *httptest.ResponseRecorder, wantBody string) {
	t.Helper()

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
