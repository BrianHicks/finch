package finch

import (
	"bytes"
	"errors"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

var (
	// ErrNoNextTask is returned when there is not a next selected task
	ErrNoNextTask = errors.New("no next task")
)

// TaskDB wraps a LevelDB instance and sets sane defaults for Finch's usage.
type TaskDB struct {
	db *leveldb.DB
	wo *opt.WriteOptions
	ro *opt.ReadOptions
}

// NewTaskDB takes a storage and returns TaskDB instance
func NewTaskDB(store storage.Storage) (*TaskDB, error) {
	tdb := new(TaskDB)

	// Open the Database with the provided Storage
	options := &opt.Options{
		Filter: filter.NewBloomFilter(15),
	}
	db, err := leveldb.Open(store, options)
	if err != nil {
		return tdb, err
	}
	tdb.db = db

	// Set default read and write options
	tdb.wo = &opt.WriteOptions{
		Sync: true,
	}
	tdb.ro = &opt.ReadOptions{
		DontFillCache: false,
	}

	return tdb, nil
}

// Close should be called on a TaskDB to end it's lifecycle. The DB should not
// be used after this is called.
func (tdb *TaskDB) Close() {
	tdb.db.Close()
	tdb.db = nil
}

// batchWriteTask makes sure that a task is completely written to the database
func (tdb *TaskDB) batchWriteTask(batch *leveldb.Batch, task *Task) error {
	szd, err := task.Serialize()
	if err != nil {
		return err
	}
	key := task.Key(TasksIndex)

	batch.Put(key.Serialize(), szd)

	for tag, presence := range task.Attrs {
		key.Index = tag
		if presence {
			batch.Put(key.Serialize(), []byte{})
		} else {
			batch.Delete(key.Serialize())
		}
	}

	return nil
}

// PutTasks inserts tasks into the database and overwrites those which
// already exist
func (tdb *TaskDB) PutTasks(tasks ...*Task) error {
	batch := new(leveldb.Batch)

	for i := 0; i < len(tasks); i++ {
		task := tasks[i]
		err := tdb.batchWriteTask(batch, task)
		if err != nil {
			return err
		}
	}

	if err := tdb.db.Write(batch, tdb.wo); err != nil {
		return err
	}

	return nil
}

// IterateOver takes an index (as prefix) to iterate over and a callback. For
// each iteration, cb will be called with the current value of the Iterator,
// and if cb returns a non-nil error that will bubble up to return from this
// function. Errors from the iterator will also be returned.
func (tdb *TaskDB) IterateOver(idx string, cb func(iterator.Iterator) error) error {
	prefix := []byte(idx)

	iter := tdb.db.NewIterator(tdb.ro)
	iter.Seek(prefix)
	defer iter.Release()

	for {
		if !bytes.HasPrefix(iter.Key(), prefix) {
			break
		}

		if err := cb(iter); err != nil {
			return err
		}

		if cont := iter.Next(); !cont {
			break
		}
	}

	if err := iter.Error(); err != nil {
		return err
	}

	return nil
}

// TasksForIndex returns a list of tasks that match an arbitrary index
// (as string)
func (tdb *TaskDB) TasksForIndex(idx string) ([]*Task, error) {
	tasks := []*Task{}
	err := tdb.IterateOver(idx, func(iter iterator.Iterator) error {
		key, err := DeserializeKey(iter.Key())
		if err != nil {
			return err
		}

		key.Index = TasksIndex

		task, err := tdb.GetTask(key)
		if err != nil {
			return err
		}
		tasks = append(tasks, task)

		return nil
	})

	return tasks, err
}

// getTaskRaw gets a task from a byteslice
func (tdb *TaskDB) getTaskRaw(key []byte) (*Task, error) {
	szd, err := tdb.db.Get(key, tdb.ro)
	if err != nil {
		return new(Task), err
	}

	task, err := DeserializeTask(szd)
	return task, err
}

// GetTask gets a single task by Key
func (tdb *TaskDB) GetTask(key *Key) (*Task, error) {
	return tdb.getTaskRaw(key.Serialize())
}

// GetPendingTasks returns a list of pending tasks
func (tdb *TaskDB) GetPendingTasks() ([]*Task, error) {
	return tdb.TasksForIndex(PendingIndex)
}

// GetSelectedTasks returns a list of currently selected tasks in
// newest-to-oldest order
func (tdb *TaskDB) GetSelectedTasks() ([]*Task, error) {
	tasks, err := tdb.TasksForIndex(SelectedIndex)
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
func (tdb *TaskDB) GetNextSelected() (*Task, error) {
	tasks, err := tdb.GetSelectedTasks()
	if err != nil {
		return new(Task), err
	}
	if len(tasks) == 0 {
		return new(Task), ErrNoNextTask
	}

	return tasks[0], nil
}
