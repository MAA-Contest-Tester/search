package main

import "gorm.io/gorm"

type StoreProblem struct {
	gorm.Model
	Statement string
	Source    string
}

type ClassificationEntry struct {
	gorm.Model
	ProblemID int
	Problem StoreProblem `gorm:"foreignKey:ProblemID;references:ID"`
	Answer int
}
