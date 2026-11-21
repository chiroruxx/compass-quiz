package main

import (
	"context"
	"errors"
	"slices"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestStudentUseCaseListStudents(t *testing.T) {
	t.Parallel()

	const facilitatorID = 1

	classA := newClass(1, "クラスA")
	students := []studentWithAccount{
		newStudentWithAccount(newStudent(1, "佐藤", classA), newStudentAccount("foo123")),
		newStudentWithAccount(newStudent(4, "加藤", classA), newStudentAccount("hoge0000")),
	}
	repositoryFailure := context.DeadlineExceeded

	cases := []struct {
		name            string
		students        []studentWithAccount
		repositoryError error
		wantCount       int
		wantError       error
	}{
		{name: "students found", students: students, repositoryError: nil, wantCount: len(students), wantError: nil},
		{name: "no students", students: nil, repositoryError: nil, wantCount: 0, wantError: errNotFound},
		{
			name: "repository error", students: nil, repositoryError: repositoryFailure,
			wantCount: 0, wantError: repositoryFailure,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			repository := NewMockStudentReader(gomock.NewController(t))
			repository.EXPECT().
				listStudentsWithAccounts(ctx, uint64(facilitatorID)).
				Return(testCase.students, testCase.repositoryError)

			got, count, err := newStudentUseCase(repository).listStudents(ctx, facilitatorID)
			if !errors.Is(err, testCase.wantError) {
				t.Fatalf("listStudents error = %v, want %v", err, testCase.wantError)
			}

			if count != testCase.wantCount || !slices.Equal(got, testCase.students) {
				t.Errorf("listStudents = (%v, %d), want (%v, %d)", got, count, testCase.students, testCase.wantCount)
			}
		})
	}
}

func TestStudentUseCaseRejectsInvalidFacilitatorID(t *testing.T) {
	t.Parallel()

	repository := NewMockStudentReader(gomock.NewController(t))
	useCase := newStudentUseCase(repository)

	_, _, err := useCase.listStudents(t.Context(), -1)
	if !errors.Is(err, errValidation) {
		t.Fatalf("listStudents error = %v, want %v", err, errValidation)
	}
}
