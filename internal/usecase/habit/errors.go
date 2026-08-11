package habit

import "errors"

var (
	ErrInvalidName        = errors.New("name is required")
	ErrInvalidExpReward   = errors.New("experience reward must be between 1 and 1000")
	ErrInvalidDailyTarget = errors.New("daily target must be greater than 0")
	ErrInvalidUnit        = errors.New("unit is required")
	ErrNotFound           = errors.New("habit not found")
	ErrForbidden          = errors.New("access denied")
	ErrInternal           = errors.New("something went wrong")
)
