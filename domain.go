package main

type student struct {
	id      int
	name    string
	account studentAccount
	class   class
}

type studentAccount struct {
	loginID string
}

type class struct {
	id   int
	name string
}
