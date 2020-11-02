package models

type User struct {
	Username string
	Password []byte
	First    string
	Last     string
	Role     string
}