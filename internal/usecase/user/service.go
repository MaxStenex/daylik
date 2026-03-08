package user

import domain "github.com/maximrozinkevich/daylik/internal/domain/user"

type service struct {
	repo domain.Repository
}

func New(repo domain.Repository) *service {
	return &service{
		repo: repo,
	}
}
