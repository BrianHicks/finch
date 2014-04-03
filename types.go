package main

import "time"

type Storage interface {
	NextID() int
	Commit() error
}

type Task struct {
	ID       int
	Desc     string
	Added    time.Time
	Delay    time.Time
	Done     bool
	Selected bool
}

type TaskStore interface {
	Storage

	SaveTask(Task) error
	GetTask(int) (*Task, error)
	AllTasks() ([]*Task, error)
	FilterTasks(func(*Task) bool) ([]*Task, error)
}
