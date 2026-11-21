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

	classA := newClass(1, "クラスA")
	students := []studentWithAccount{
		newStudentWithAccount(newStudent(1, "佐藤", classA), newStudentAccount("foo123")),
		newStudentWithAccount(newStudent(4, "加藤", classA), newStudentAccount("hoge0000")),
	}
	repositoryFailure := context.DeadlineExceeded

	cases := []struct {
		name    string
		arrange func(*testing.T) (studentUseCase, int)
		assert  func(*testing.T, []studentWithAccount, int, error)
	}{
		{
			name:    "students found",
			arrange: arrangeStudentLookup(students, nil),
			assert:  assertListStudentsResult(students, nil),
		},
		{
			name:    "no students",
			arrange: arrangeStudentLookup(nil, nil),
			assert:  assertListStudentsResult(nil, errNotFound),
		},
		{
			name:    "repository error",
			arrange: arrangeStudentLookup(nil, repositoryFailure),
			assert:  assertListStudentsResult(nil, repositoryFailure),
		},
		{
			name: "invalid facilitator ID",
			arrange: func(t *testing.T) (studentUseCase, int) {
				t.Helper()

				repository := NewMockStudentReader(gomock.NewController(t))

				return newStudentUseCase(repository), -1
			},
			assert: assertListStudentsResult(nil, errValidation),
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			useCase, facilitatorID := testCase.arrange(t)

			got, count, err := useCase.listStudents(t.Context(), facilitatorID)

			testCase.assert(t, got, count, err)
		})
	}
}

func arrangeStudentLookup(students []studentWithAccount, err error) func(*testing.T) (studentUseCase, int) {
	return func(t *testing.T) (studentUseCase, int) {
		t.Helper()

		const facilitatorID = 1

		repository := NewMockStudentReader(gomock.NewController(t))
		repository.EXPECT().listStudentsWithAccounts(t.Context(), uint64(facilitatorID)).Return(students, err)

		return newStudentUseCase(repository), facilitatorID
	}
}

func assertListStudentsResult(
	want []studentWithAccount, wantError error,
) func(*testing.T, []studentWithAccount, int, error) {
	return func(t *testing.T, got []studentWithAccount, count int, err error) {
		t.Helper()

		if !errors.Is(err, wantError) {
			t.Fatalf("listStudents error = %v, want %v", err, wantError)
		}

		if count != len(want) || !slices.Equal(got, want) {
			t.Errorf("listStudents = (%v, %d), want (%v, %d)", got, count, want, len(want))
		}
	}
}
