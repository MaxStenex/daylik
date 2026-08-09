package habit_log

import "errors"

var (
	ErrInvalidHabit           = errors.New("Please select a valid habit")
	ErrInvalidCompletedCount  = errors.New("Completed count must be greater than 0")
	ErrCompletedCountTooLarge = errors.New("Completed count cannot exceed the daily target")
	ErrNotFound               = errors.New("Habit log not found")
	ErrForbidden              = errors.New("Access denied")
	ErrUpdateOutdated         = errors.New("You can only update today's log")
	ErrInternal               = errors.New("Something went wrong")
)
