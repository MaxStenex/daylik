package habit_log

import "errors"

var (
	ErrInvalidHabit           = errors.New("please select a valid habit")
	ErrInvalidCompletedCount  = errors.New("completed count must be greater than 0")
	ErrCompletedCountTooLarge = errors.New("completed count cannot exceed the daily target")
	ErrNotFound               = errors.New("habit log not found")
	ErrForbidden              = errors.New("access denied")
	ErrUpdateOutdated         = errors.New("you can only update today's log")
	ErrInternal               = errors.New("something went wrong")
)
