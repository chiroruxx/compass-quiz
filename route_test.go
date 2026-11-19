package main

import (
	"database/sql"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/go-sql-driver/mysql" // Register the MySQL database/sql driver.
)

func TestHealthRoute(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("requires the Compose MySQL test database")
	}

	insertTestFixture(t)

	server := httptest.NewServer(newRouter())
	t.Cleanup(server.Close)

	request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, server.URL+"/healthz", nil)
	if err != nil {
		t.Fatalf("create health request: %v", err)
	}

	response, err := server.Client().Do(request)
	if err != nil {
		t.Fatalf("send health request: %v", err)
	}

	t.Cleanup(func() {
		err := response.Body.Close()
		if err != nil {
			t.Errorf("close response body: %v", err)
		}
	})

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("read health response: %v", err)
	}

	assertHealthz(t, response.StatusCode, response.Header.Get("Content-Type"), string(body))
}

func insertTestFixture(t *testing.T) {
	t.Helper()

	const testDSN = "compass:compass_local@tcp(mysql:3306)/compass_quiz_test"

	database, err := sql.Open("mysql", testDSN)
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}

	t.Cleanup(func() {
		err := database.Close()
		if err != nil {
			t.Errorf("close test database: %v", err)
		}
	})

	transaction, err := database.BeginTx(t.Context(), nil)
	if err != nil {
		t.Fatalf("start fixture transaction: %v", err)
	}

	t.Cleanup(func() {
		// database/sql may already have rolled back when the test context ended.
		err := transaction.Rollback()
		if err != nil && !errors.Is(err, sql.ErrTxDone) {
			t.Errorf("roll back test fixture: %v", err)
		}
	})

	_, err = transaction.ExecContext(t.Context(), "INSERT INTO classes (id, name) VALUES (9000, 'Test class')")
	if err != nil {
		t.Fatalf("insert test fixture: %v", err)
	}
}
