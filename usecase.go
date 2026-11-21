package main

import (
	"context"
	"fmt"
)

//go:generate go tool -modfile=tools/mockgen/go.mod mockgen -source=usecase.go -destination=usecase_mock_test.go -package=main -mock_names=studentReader=MockStudentReader -write_package_comment=false

type studentReader interface {
	listStudentsWithAccounts(ctx context.Context, facilitatorID uint64) ([]studentWithAccount, error)
}

type studentUseCase struct {
	students studentReader
}

func newStudentUseCase(students studentReader) studentUseCase {
	return studentUseCase{students: students}
}

func (u studentUseCase) listStudents(ctx context.Context, facilitatorID int) ([]studentWithAccount, int, error) {
	if facilitatorID < 1 {
		return nil, 0, fmt.Errorf("%w: facilitator ID must be at least 1", errValidation)
	}

	students, err := u.students.listStudentsWithAccounts(ctx, uint64(facilitatorID))
	if err != nil {
		return nil, 0, err
	}

	if len(students) == 0 {
		return nil, 0, fmt.Errorf("students: %w", errNotFound)
	}

	return students, len(students), nil
}
