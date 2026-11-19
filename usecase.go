package main

type useCase struct{}

func newUseCase() useCase {
	return useCase{}
}

// listStudents returns the assignment example until the MySQL lookup is implemented.
func (useCase) listStudents(_ int) ([]studentWithAccount, int) {
	students := []studentWithAccount{
		newStudentWithAccount(
			newStudent(1, "佐藤", newClass(1, "クラスA")),
			newStudentAccount("foo123"),
		),
	}

	return students, len(students)
}
