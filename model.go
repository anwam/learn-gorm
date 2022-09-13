package main

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string
	Username string
	Password string `json:"-"`
}

type Book struct {
	gorm.Model
	Name     string
	AuthorID uint
	Author   User
}
