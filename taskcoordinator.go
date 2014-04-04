package main

import "time"

type TaskCoordinator struct {
	storage TaskStore
}

func (tc *TaskCoordinator) Add(desc string) (*Task, error) {
	t := &Task{
		Desc:   desc,
		Active: time.Now(),
	}

	err := tc.storage.SaveTask(t)
	return t, err
}
