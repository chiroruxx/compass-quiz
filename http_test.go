package main

import (
	"database/sql"
	"encoding/json/jsontext"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/go-testfixtures/testfixtures/v3"
)

const wantAssignedStudentsJSON = `{
	"students": [
		{"id": 1, "name": "佐藤", "loginId": "foo123", "classroom": {"id": 1, "name": "クラスA"}},
		{"id": 4, "name": "加藤", "loginId": "hoge0000", "classroom": {"id": 1, "name": "クラスA"}},
		{"id": 7, "name": "小田原", "loginId": "fizz9999", "classroom": {"id": 1, "name": "クラスA"}},
		{"id": 10, "name": "Williams", "loginId": "xxx:42", "classroom": {"id": 1, "name": "クラスA"}}
	],
	"totalCount": 4
}`

//nolint:paralleltest // Seeded cases modify the same tables and must run sequentially.
func TestGetStudentsHTTP(t *testing.T) {
	assertBadRequest := assertEmptyHTTPResponse(http.StatusBadRequest)

	cases := []struct {
		name    string
		path    string
		arrange func(*testing.T, *sql.DB)
		assert  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "assigned students", path: "/students?facilitator_id=1",
			arrange: loadStudentFixtures,
			assert: func(t *testing.T, response *httptest.ResponseRecorder) {
				t.Helper()

				assertHTTPStatus(t, response, http.StatusOK)
				assertStudentsJSON(t, response, wantAssignedStudentsJSON)
			},
		},
		{
			name: "unknown facilitator", path: "/students?facilitator_id=999",
			arrange: loadStudentFixtures, assert: assertEmptyHTTPResponse(http.StatusNotFound),
		},
		{
			name: "internal server error", path: "/students?facilitator_id=1",
			arrange: closeHTTPTestDB, assert: assertEmptyHTTPResponse(http.StatusInternalServerError),
		},
		{name: "missing", path: "/students", arrange: nil, assert: assertBadRequest},
		{name: "empty", path: "/students?facilitator_id=", arrange: nil, assert: assertBadRequest},
		{name: "text", path: "/students?facilitator_id=abc", arrange: nil, assert: assertBadRequest},
		{name: "fraction", path: "/students?facilitator_id=1.5", arrange: nil, assert: assertBadRequest},
		{name: "zero", path: "/students?facilitator_id=0", arrange: nil, assert: assertBadRequest},
		{name: "negative", path: "/students?facilitator_id=-1", arrange: nil, assert: assertBadRequest},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			pool := openHTTPTestDB(t)
			if testCase.arrange != nil {
				testCase.arrange(t, pool)
			}

			router := newRouter(newStudentUseCase(newStudentRepository(pool)))
			request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, testCase.path, nil)
			response := httptest.NewRecorder()

			router.ServeHTTP(response, request)

			testCase.assert(t, response)
		})
	}
}

func openHTTPTestDB(t *testing.T) *sql.DB {
	t.Helper()

	config := mysqlConfigFromEnv()
	config.DBName = "compass_quiz_test"

	pool := sql.OpenDB(txdb.New("mysql", config.FormatDSN()))

	t.Cleanup(func() {
		closeHTTPTestDB(t, pool)
	})

	return pool
}

func closeHTTPTestDB(t *testing.T, pool *sql.DB) {
	t.Helper()

	err := pool.Close()
	if err != nil {
		t.Errorf("close test database: %v", err)
	}
}

func loadStudentFixtures(t *testing.T, pool *sql.DB) {
	t.Helper()

	fixtures, err := testfixtures.New(
		testfixtures.Database(pool),
		testfixtures.Dialect("mysql"),
		testfixtures.FilesMultiTables("testdata/GetStudents.yml"),
		// MySQL resets sequences with ALTER TABLE, which implicitly commits the transaction.
		testfixtures.SkipResetSequences(),
		testfixtures.SkipTableChecksumComputation(),
	)
	if err != nil {
		t.Fatalf("create fixture loader: %v", err)
	}

	err = fixtures.Load()
	if err != nil {
		t.Fatalf("load student fixtures: %v", err)
	}
}

func assertHTTPStatus(t *testing.T, response *httptest.ResponseRecorder, wantStatus int) {
	t.Helper()

	if response.Code != wantStatus {
		t.Fatalf("status = %d, want %d", response.Code, wantStatus)
	}
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

func assertEmptyHTTPResponse(wantStatus int) func(*testing.T, *httptest.ResponseRecorder) {
	return func(t *testing.T, response *httptest.ResponseRecorder) {
		t.Helper()

		assertHTTPStatus(t, response, wantStatus)

		if response.Body.Len() != 0 {
			t.Errorf("body = %q, want empty", response.Body.String())
		}
	}
}
