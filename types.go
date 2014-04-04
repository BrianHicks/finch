package main

import "errors"

var (
	NoSuchTask = errors.New("no such task")
)

type Storage interface {
	NextID() uint
	Commit() error
}

type TaskStore interface {
	Storage

	SaveTask(*Task) error
	GetTask(string) (*Task, error)
	FilterTasks(func(*Task) bool) ([]*Task, error)
}
