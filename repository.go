package main

import (
	"context"
	"fmt"

	"compass-quiz/db"
)

type studentRepository struct {
	queries *db.Queries
}

func newStudentRepository(conn db.DBTX) studentRepository {
	return studentRepository{queries: db.New(conn)}
}

func (r studentRepository) listStudentsWithAccounts(
	ctx context.Context, facilitatorID uint64,
) ([]studentWithAccount, error) {
	rows, err := r.queries.ListStudentsWithAccounts(ctx, facilitatorID)
	if err != nil {
		return nil, fmt.Errorf("list students by facilitator: %w", err)
	}

	students := make([]studentWithAccount, len(rows))
	for i, row := range rows {
		students[i] = newStudentWithAccount(
			newStudent(row.ID, row.Name, newClass(row.ClassID, row.ClassName)),
			newStudentAccount(row.LoginID),
		)
	}

	return students, nil
}
