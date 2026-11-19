package main

type student struct {
	id      int
	name    string
	account studentAccount
	class   class
}

func newStudent(id int, name string, account studentAccount, class class) student {
	return student{id: id, name: name, account: account, class: class}
}

type studentAccount struct {
	loginID string
}

func newStudentAccount(loginID string) studentAccount {
	return studentAccount{loginID: loginID}
}

type class struct {
	id   int
	name string
}

func newClass(id int, name string) class {
	return class{id: id, name: name}
}
