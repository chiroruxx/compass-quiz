package main

type useCase struct{}

func newUseCase() useCase {
	return useCase{}
}

// listStudents returns the assignment example until the MySQL lookup is implemented.
func (useCase) listStudents(_ int) ([]student, int) {
	students := []student{
		newStudent(
			1,
			"佐藤", //nolint:gosmopolitan // Use the name from the assignment's response example.
			newStudentAccount("foo123"),
			newClass(1, "クラスA"),
		),
	}

	return students, len(students)
}
