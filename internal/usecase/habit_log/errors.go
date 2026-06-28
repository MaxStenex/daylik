package habit_log

import "errors"

var (
	ErrInvalidHabit           = errors.New("invalid habit id")
	ErrInvalidCompletedCount  = errors.New("completed count must be greater than 0")
	ErrCompletedCountTooLarge = errors.New("completed count cannot be greater than daily target")
	ErrNotFound               = errors.New("habit log not found")
	ErrForbidden              = errors.New("access denied")
	ErrUpdateOutdated         = errors.New("cannot update logs from previous days")
)
