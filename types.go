package main

import "errors"

var (
	// ErrNoSuchTask is returned when a task cannot be found
	ErrNoSuchTask = errors.New("no such task")

	// ErrNoSuchKey is returned when a key cannot be found
	ErrNoSuchKey = errors.New("no such key")
)

// Storage is the base interface
type Storage interface {
	NextID() uint
	Commit() error
}

// TaskStorage deals with the storage of tasks. Wow, revolutionary!
type TaskStorage interface {
	SaveTask(...*Task) error
	GetTask(string) (*Task, error)
	DeleteTask(...string) error
	FilterTasks(func(*Task) bool) ([]*Task, error)
}

// MetaStore deals with the storage of meta-information. You might want to
// store sync data in this, for example. It's all strings, though.
type MetaStore interface {
	SetMeta(string, string) error
	GetMeta(string) (string, error)
}

// TaskStore is a Storage + a TaskStorage
type TaskStore interface {
	Storage
	TaskStorage
}

// MetaTaskStore is a Storage + TaskStorage + MetaStore
type MetaTaskStore interface {
	Storage

	TaskStorage
	MetaStore
}
