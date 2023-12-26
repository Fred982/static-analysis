package models

import "time"

type Tabler interface {
	TableName() string
}

type GithubRepository struct {
	ID          int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
	FullName    string
	Description string
}

func (GithubRepository) TableName() string {
	return "github_repository"
}
