package habit

import (
	domain "github.com/maximrozinkevich/daylik/internal/domain/habit"
)

type service struct {
	habitRepo domain.Repository
}

func New(habitRepo domain.Repository) *service {
	return &service{habitRepo: habitRepo}
}
