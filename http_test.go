package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
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
	loadStudentFixtures := func(t *testing.T, pool *sql.DB) {
		t.Helper()

		loadHTTPFixtures(t, pool, "testdata/GetStudents.yml")
	}
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
				assertJSONResponse(t, response, wantAssignedStudentsJSON)
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
