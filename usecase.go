package main

type useCase struct{}

// listStudents returns the assignment example until the MySQL lookup is implemented.
func (useCase) listStudents(_ int) ([]student, int) {
	students := []student{
		{
			id:   1,
			name: "佐藤", //nolint:gosmopolitan // Use the name from the assignment's response example.
			account: studentAccount{
				loginID: "foo123",
			},
			class: class{
				id:   1,
				name: "クラスA",
			},
		},
	}

	return students, len(students)
}
