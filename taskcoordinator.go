package main

import (
	"sort"
	"time"
)

type TaskCoordinator struct {
	storage TaskStore
}

func (tc *TaskCoordinator) Close() error {
	return tc.storage.Commit()
}

func (tc *TaskCoordinator) Add(desc string) *Task {
	t := &Task{
		Desc:   desc,
		Active: time.Now(),
	}

	return t
}

func (tc *TaskCoordinator) Get(id string) (*Task, error) {
	return tc.storage.GetTask(id)
}

func (tc *TaskCoordinator) Delay(id string, until time.Time) error {
	t, err := tc.storage.GetTask(id)
	if err != nil {
		return err
	}

	t.Active = until
	t.Selected = false
	err = tc.storage.SaveTask(t)
	if err != nil {
		return err
	}

	return nil
}

func (tc *TaskCoordinator) Select(ids ...string) error {
	// make sure we have all the tasks before we perform the operation
	tasks := []*Task{}
	for _, id := range ids {
		t, err := tc.storage.GetTask(id)
		if err != nil {
			return err
		}

		tasks = append(tasks, t)
	}

	// select all the tasks and save them
	for _, t := range tasks {
		t.Selected = true
	}

	err := tc.storage.SaveTask(tasks...)
	if err != nil {
		return err
	}

	return nil
}

func (tc *TaskCoordinator) MarkDone(ids ...string) error {
	// make sure we have all the tasks before we perform the operation
	tasks := []*Task{}
	for _, id := range ids {
		t, err := tc.storage.GetTask(id)
		if err != nil {
			return err
		}

		tasks = append(tasks, t)
	}

	// mark all the tasks done and save them
	for _, t := range tasks {
		t.MarkDone()
	}

	err := tc.storage.SaveTask(tasks...)
	if err != nil {
		return err
	}

	return nil
}

func (tc *TaskCoordinator) Save(tasks ...*Task) error {
	for _, t := range tasks {
		err := tc.storage.SaveTask(t)
		if err != nil {
			return err
		}
	}

	return nil
}

func (tc *TaskCoordinator) Delete(ids ...string) error {
	return tc.storage.DeleteTask(ids...)
}

func (tc *TaskCoordinator) Selected() ([]*Task, error) {
	tasks, err := tc.storage.FilterTasks(func(t *Task) bool { return t.Selected })

	if len(tasks) == 0 {
		return tasks, NoSuchTask
	}

	sort.Sort(sort.Reverse(ByActive(tasks)))

	return tasks, err
}

func (tc *TaskCoordinator) NextSelected() (*Task, error) {
	tasks, err := tc.Selected()
	if err != nil {
		return nil, err
	}

	return tasks[0], nil
}

func (tc *TaskCoordinator) Available() ([]*Task, error) {
	now := time.Now()
	tasks, err := tc.storage.FilterTasks(func(t *Task) bool { return !t.Done && now.After(t.Active) })
	if err != nil {
		return tasks, err
	}

	if len(tasks) == 0 {
		return tasks, NoSuchTask
	}

	sort.Sort(ByActive(tasks))

	return tasks, nil
}
