package main

import "errors"

var (
	NoSuchTask = errors.New("no such task")
	NoSuchKey  = errors.New("no such key")
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

type MetaStore interface {
	Storage

	SetMeta(string, string) error
	GetMeta(string) (string, error)
}
