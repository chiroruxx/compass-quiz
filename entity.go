package main

type student struct {
	id    uint64
	name  string
	class class
}

func newStudent(id uint64, name string, class class) student {
	return student{id: id, name: name, class: class}
}

type studentWithAccount struct {
	student

	account studentAccount
}

func newStudentWithAccount(student student, account studentAccount) studentWithAccount {
	return studentWithAccount{student: student, account: account}
}

type studentAccount struct {
	loginID string
}

func newStudentAccount(loginID string) studentAccount {
	return studentAccount{loginID: loginID}
}

type class struct {
	id   uint64
	name string
}

func newClass(id uint64, name string) class {
	return class{id: id, name: name}
}
