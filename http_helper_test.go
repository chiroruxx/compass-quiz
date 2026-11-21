package main

import (
	"database/sql"
	"encoding/json/jsontext"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/go-testfixtures/testfixtures/v3"
)

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

func loadHTTPFixtures(t *testing.T, pool *sql.DB, fixturePath string) {
	t.Helper()

	fixtures, err := testfixtures.New(
		testfixtures.Database(pool),
		testfixtures.Dialect("mysql"),
		testfixtures.FilesMultiTables(fixturePath),
		// MySQL resets sequences with ALTER TABLE, which implicitly commits the transaction.
		testfixtures.SkipResetSequences(),
		testfixtures.SkipTableChecksumComputation(),
	)
	if err != nil {
		t.Fatalf("create fixture loader for %q: %v", fixturePath, err)
	}

	err = fixtures.Load()
	if err != nil {
		t.Fatalf("load fixtures from %q: %v", fixturePath, err)
	}
}

func assertHTTPStatus(t *testing.T, response *httptest.ResponseRecorder, wantStatus int) {
	t.Helper()

	if response.Code != wantStatus {
		t.Fatalf("status = %d, want %d", response.Code, wantStatus)
	}
}

func assertJSONResponse(t *testing.T, response *httptest.ResponseRecorder, wantBody string) {
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
