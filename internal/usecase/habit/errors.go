package habit

import "errors"

var (
	ErrInvalidName        = errors.New("Name is required")
	ErrInvalidExpReward   = errors.New("Experience reward must be between 1 and 1000")
	ErrInvalidDailyTarget = errors.New("Daily target must be greater than 0")
	ErrInvalidUnit        = errors.New("Unit is required")
	ErrNotFound           = errors.New("Habit not found")
	ErrForbidden          = errors.New("Access denied")
	ErrInternal           = errors.New("Something went wrong")
)
