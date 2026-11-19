//go:build integration

package main_test

import (
	"database/sql"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	api "compass-quiz"

	_ "github.com/go-sql-driver/mysql" // Register the MySQL database/sql driver.
)

func TestIntegrationEnvironment(t *testing.T) {
	t.Parallel()
	checkTestDatabase(t)

	server := httptest.NewServer(api.NewRouter())
	t.Cleanup(server.Close)

	request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, server.URL+healthPath, nil)
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

func checkTestDatabase(t *testing.T) {
	t.Helper()

	dsn := os.Getenv("TEST_MYSQL_DSN")
	if dsn == "" {
		t.Fatal("TEST_MYSQL_DSN must be set to the connection string for compass_quiz_test")
	}

	database, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}

	t.Cleanup(func() {
		err := database.Close()
		if err != nil {
			t.Errorf("close test database: %v", err)
		}
	})

	var databaseName string

	err = database.QueryRowContext(t.Context(), "SELECT DATABASE()").Scan(&databaseName)
	if err != nil {
		t.Fatalf("identify test database: %v", err)
	}

	if databaseName != "compass_quiz_test" {
		t.Fatalf("connected to %q, want compass_quiz_test", databaseName)
	}

	var seeded bool

	err = database.QueryRowContext(t.Context(), "SELECT EXISTS (SELECT 1 FROM students)").Scan(&seeded)
	if err != nil {
		t.Fatalf("read student seed data: %v", err)
	}

	if !seeded {
		t.Fatal("the test database contains no student seed data")
	}
}
