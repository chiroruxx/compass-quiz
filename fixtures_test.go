package main

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/go-sql-driver/mysql"
	"github.com/go-testfixtures/testfixtures/v3"
)

func openHTTPTestDB(t *testing.T) *sql.DB {
	t.Helper()

	config := testMySQLConfig()

	pool := sql.OpenDB(txdb.New("mysql", config.FormatDSN()))

	t.Cleanup(func() {
		closeHTTPTestDB(t, pool)
	})

	return pool
}

func testMySQLConfig() *mysql.Config {
	config := mysqlConfigFromEnv()
	config.DBName = "compass_quiz_test"

	return config
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

//nolint:paralleltest // HTTP tests load the same fixture rows.
func TestStudentFixturesRollback(t *testing.T) {
	pool, err := openMySQL(testMySQLConfig())
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}

	t.Cleanup(func() {
		closeHTTPTestDB(t, pool)
	})

	before := countFixtureStudents(t, pool)
	t.Run("load fixtures", func(t *testing.T) {
		transaction := openHTTPTestDB(t)
		loadStudentFixtures(t, transaction)

		const wantStudents = 5
		if count := countFixtureStudents(t, transaction); count != wantStudents {
			t.Fatalf("students in the transaction = %d, want %d", count, wantStudents)
		}
	})

	if after := countFixtureStudents(t, pool); after != before {
		t.Errorf("students after rollback = %d, want %d", after, before)
	}
}

func countFixtureStudents(t *testing.T, pool *sql.DB) int {
	t.Helper()

	var count int

	err := pool.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM students").Scan(&count)
	if err != nil {
		t.Fatalf("count fixture students: %v", err)
	}

	return count
}
