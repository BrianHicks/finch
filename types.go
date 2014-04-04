package main

import (
	"errors"
	"time"
)

var (
	NoSuchTask = errors.New("no such task")
)

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
	GetTask(string) (*Task, error)
	FilterTasks(func(*Task) bool) ([]*Task, error)
}
