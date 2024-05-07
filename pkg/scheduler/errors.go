package scheduler

import "errors"

var (
	ErrInvalidTask    = errors.New("task is invalid")
	ErrTaskNotFound   = errors.New("task not found")
	ErrWorkerNotFound = errors.New("worker not found")
)
