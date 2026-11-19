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
			"佐藤",
			newStudentAccount("foo123"),
			newClass(1, "クラスA"),
		),
	}

	return students, len(students)
}
