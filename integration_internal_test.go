//go:build integration

package main

import (
	"database/sql"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/go-sql-driver/mysql" // Register the MySQL database/sql driver.
)

func TestIntegrationEnvironment(t *testing.T) {
	t.Parallel()
	checkTestDatabase(t)

	server := httptest.NewServer(newRouter())
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

	var seeded bool

	err = database.QueryRowContext(t.Context(), "SELECT EXISTS (SELECT 1 FROM students)").Scan(&seeded)
	if err != nil {
		t.Fatalf("read student seed data: %v", err)
	}

	if !seeded {
		t.Fatal("the test database contains no student seed data")
	}
}
