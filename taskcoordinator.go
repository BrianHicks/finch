package main

import "time"

type TaskCoordinator struct {
	storage MetaTaskStore
}

func (tc *TaskCoordinator) Add(desc string) (*Task, error) {
	t := &Task{
		Desc:   desc,
		Active: time.Now(),
	}

	err := tc.storage.SaveTask(t)
	return t, err
}

func (tc *TaskCoordinator) Delay(id string, until time.Time) error {
	t, err := tc.storage.GetTask(id)
	if err != nil {
		return err
	}

	t.Active = until
	err = tc.storage.SaveTask(t)
	if err != nil {
		return err
	}

	return nil
}
