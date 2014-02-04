package core

import (
	"errors"

	"github.com/BrianHicks/finch/persist"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

var (
	// ErrNoTask is returned when there is not a next selected task
	ErrNoTask = errors.New("no such task")
)

// TaskStore wraps a LevelStore instance and sets sane defaults for Finch's usage.
type TaskStore struct {
	Store *persist.Store
	wo    *opt.WriteOptions
	ro    *opt.ReadOptions
}

// NewTaskStore takes a storage and returns TaskStore instance
func NewTaskStore(storage storage.Storage) (*TaskStore, error) {
	ts := new(TaskStore)

	store, err := persist.New(storage)
	if err != nil {
		return ts, err
	}
	ts.Store = store

	return ts, nil
}

// Close should be called on a TaskStore to end it's lifecycle. The Store should not
// be used after this is called.
func (ts *TaskStore) Close() {
	ts.Store.Close()
	ts.Store = nil
}

// batchWriteTask makes sure that a task is completely written to the database
func (ts *TaskStore) batchWriteTask(batch *persist.LoggedBatch, task *Task) error {
	szd, err := task.Serialize()
	if err != nil {
		return err
	}
	key := task.Key()

	batch.Put(key.Serialize(TasksIndex), szd)

	for tag, presence := range task.Attrs {
		if presence {
			batch.Put(key.Serialize(tag), []byte{})
		} else {
			batch.Delete(key.Serialize(tag))
		}
	}

	return nil
}

// PutTasks inserts tasks into the database and overwrites those which
// already exist
func (ts *TaskStore) PutTasks(tasks ...*Task) error {
	batch := persist.NewLoggedBatch()

	for i := 0; i < len(tasks); i++ {
		task := tasks[i]
		err := ts.batchWriteTask(batch, task)
		if err != nil {
			return err
		}
	}

	if err := ts.Store.Write(batch.Batch); err != nil {
		return err
	}

	return nil
}

// MoveTask reindexes the current Task if the components of the Key change. It
// does this in a Batch, so the Move is an atomic operation.
//
// Currently, that means if you change Task.Timestamp or Task.ID you need to use
// this or old data will always show up.
func (ts *TaskStore) MoveTask(oldKey *Key, task *Task) error {
	batch := persist.NewLoggedBatch()

	batch.Delete(oldKey.Serialize(TasksIndex))
	for prefix := range task.Attrs {
		batch.Delete(oldKey.Serialize(prefix))
	}

	err := ts.batchWriteTask(batch, task)
	if err != nil {
		return err
	}

	if err := ts.Store.Write(batch.Batch); err != nil {
		return err
	}

	return nil
}

// TasksForIndex returns a list of tasks that match an arbitrary index
// (as string)
func (ts *TaskStore) TasksForIndex(prefix string) ([]*Task, error) {
	tasks := []*Task{}

	raw, err := ts.Store.Prefix([]byte(prefix)).All()
	if err != nil {
		return tasks, err
	}

	for i := 0; i < len(raw); i++ {
		key, err := DeserializeKey(raw[i].Key)
		if err != nil {
			return tasks, err
		}

		task, err := ts.GetTask(key)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, task)
	}

	return tasks, err
}

// getTaskRaw gets a task from a byteslice
func (ts *TaskStore) getTaskRaw(key []byte) (*Task, error) {
	szd, err := ts.Store.Get(key)
	if len(szd) == 0 {
		return new(Task), ErrNoTask
	}
	if err != nil {
		return new(Task), err
	}

	task, err := DeserializeTask(szd)
	return task, err
}

// GetTask gets a single task by Key
func (ts *TaskStore) GetTask(key *Key) (*Task, error) {
	return ts.getTaskRaw(key.Serialize(TasksIndex))
}

// GetPendingTasks returns a list of pending tasks
func (ts *TaskStore) GetPendingTasks() ([]*Task, error) {
	return ts.TasksForIndex(TagPending)
}

// GetSelectedTasks returns a list of currently selected tasks in
// newest-to-oldest order
func (ts *TaskStore) GetSelectedTasks() ([]*Task, error) {
	tasks, err := ts.TasksForIndex(TagSelected)
	if err != nil {
		return tasks, err
	}

	reversed := []*Task{}
	for i := len(tasks) - 1; i >= 0; i-- {
		reversed = append(reversed, tasks[i])
	}
	return reversed, nil
}

// GetNextSelected returns the next (most recent) selected Task
func (ts *TaskStore) GetNextSelected() (*Task, error) {
	tasks, err := ts.GetSelectedTasks()
	if err != nil {
		return new(Task), err
	}
	if len(tasks) == 0 {
		return new(Task), ErrNoTask
	}

	return tasks[0], nil
}
