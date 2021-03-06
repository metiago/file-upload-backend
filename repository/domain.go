package repository

import (
	"time"
)

// Privilege is a type to define permissions eg.: READ, WRITE and DELETE
type Privilege struct {
	ID      int       `json:"id"`
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
}

// File is a type to define a stored file
type File struct {
	ID         int          `json:"id"`
	Name       string       `json:"name"`
	Ext        string       `json:"ext"`
	Privileges *[]Privilege `json:"permissions"`
	Data       []byte       `json:"-"`
	Created    time.Time    `json:"created"`
	Username   string       `json:"username"`
}

// User is a type to define users
type User struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	Username        string    `json:"username"`
	Password        string    `json:"password"`
	ConfirmPassword string    `json:"confirm_password"`
	UpdatedPassword string    `json:"updated_password"`
	Created         time.Time `json:"created"`
	Files           []*File   `json:"files"`
}
