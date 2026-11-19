package main

type student struct {
	id      uint64
	name    string
	account studentAccount
	class   class
}

func newStudent(id uint64, name string, account studentAccount, class class) student {
	return student{id: id, name: name, account: account, class: class}
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
