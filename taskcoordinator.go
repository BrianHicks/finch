package main

import (
	"sort"
	"time"
)

// TaskCoordinator is the layer between Tasks and TaskStores. It coordinates
// operations on multiple tasks and implements common tasks.
type TaskCoordinator struct {
	storage TaskStore
}

// Close should be called at the end of TaskCoordinator's lifecycle.
func (tc *TaskCoordinator) Close() error {
	return tc.storage.Commit()
}

// Add creates a new Task and returns it
func (tc *TaskCoordinator) Add(desc string) *Task {
	t := &Task{
		Desc:   desc,
		Active: time.Now(),
	}

	return t
}

// Get gets a single task by ID
func (tc *TaskCoordinator) Get(id string) (*Task, error) {
	return tc.storage.GetTask(id)
}

// Delay delays a task until some later date (or until now, if you just want to
// take it off of your "current" list)
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

// Select takes multiple IDs and marks them as selected
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

// MarkDone takes multiple IDs and marks them all as done.
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

// Save saves multiple tasks to the underlying store
func (tc *TaskCoordinator) Save(tasks ...*Task) error {
	for _, t := range tasks {
		err := tc.storage.SaveTask(t)
		if err != nil {
			return err
		}
	}

	return nil
}

// Delete taskes multiple IDs to delete
func (tc *TaskCoordinator) Delete(ids ...string) error {
	return tc.storage.DeleteTask(ids...)
}

// Selected gets all currently selected tasks, in reverse active order (as
// according to Final Version)
func (tc *TaskCoordinator) Selected() ([]*Task, error) {
	tasks, err := tc.storage.FilterTasks(func(t *Task) bool { return t.Selected })

	if len(tasks) == 0 {
		return tasks, ErrNoSuchTask
	}

	sort.Sort(sort.Reverse(ByActive(tasks)))

	return tasks, err
}

// NextSelected gets the first task from `Selected`.
func (tc *TaskCoordinator) NextSelected() (*Task, error) {
	tasks, err := tc.Selected()
	if err != nil {
		return nil, err
	}

	return tasks[0], nil
}

// Available returns all available tasks. That is to say, ones that are not
// active sometime in the future and are not done.
func (tc *TaskCoordinator) Available() ([]*Task, error) {
	now := time.Now()
	tasks, err := tc.storage.FilterTasks(func(t *Task) bool { return !t.Done && now.After(t.Active) })
	if err != nil {
		return tasks, err
	}

	if len(tasks) == 0 {
		return tasks, ErrNoSuchTask
	}

	sort.Sort(ByActive(tasks))

	return tasks, nil
}
