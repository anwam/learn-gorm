package main

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string
	Username string
	Password string `json:"-"`
	Books    []Book `gorm:"foreignKey:AuthorID"`
}

type Book struct {
	gorm.Model
	Name     string
	AuthorID uint
}
