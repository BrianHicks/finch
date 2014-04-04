package main

import "time"

type Storage interface {
	NextID() uint
	Commit() error
}

type Task struct {
	ID       string
	Desc     string
	Active   time.Time
	Done     bool
	Selected bool
}

type TaskStore interface {
	Storage

	SaveTask(*Task) error
	GetTask(int) (*Task, error)
	AllTasks() ([]*Task, error)
	FilterTasks(func(*Task) bool) ([]*Task, error)
}
