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

type TaskStorage interface {
	SaveTask(...*Task) error
	GetTask(string) (*Task, error)
	DeleteTask(...string) error
	FilterTasks(func(*Task) bool) ([]*Task, error)
}

type MetaStore interface {
	SetMeta(string, string) error
	GetMeta(string) (string, error)
}

type TaskStore interface {
	Storage
	TaskStorage
}

type MetaTaskStore interface {
	Storage

	TaskStorage
	MetaStore
}
