package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"sync"

	"github.com/cryptix/goremutake"
)

var errNoSuch = errors.New("no such file or directory")

// JSONStore is a Storage implementation on a disk-backed JSON file
type JSONStore struct {
	filename string

	CurID uint `json:"cur_id"`

	Tasks    map[string]*Task `json:"tasks"`
	taskLock *sync.RWMutex

	Meta     map[string]string `json:"meta"`
	metaLock *sync.RWMutex
}

// NewJSONStore properly initializes a JSONStore
func NewJSONStore(filename string) (*JSONStore, error) {
	store := &JSONStore{
		filename: filename,
		taskLock: new(sync.RWMutex),
		metaLock: new(sync.RWMutex),
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		if perr, ok := err.(*os.PathError); ok {
			if perr.Err.Error() == errNoSuch.Error() {
				store.Tasks = map[string]*Task{}
				store.Meta = map[string]string{}
				return store, nil
			}
		}
		return store, err
	}

	err = json.Unmarshal(content, store)
	if err != nil {
		return store, err
	}

	return store, nil
}

// NextID gets the next ID in sequence for the JSON store
func (j *JSONStore) NextID() uint {
	j.CurID++
	return j.CurID
}

// Commit writes this data to the disk location for which it is configured.
func (j *JSONStore) Commit() error {
	j.taskLock.Lock()
	defer j.taskLock.Unlock()

	bytes, err := json.Marshal(j)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(j.filename, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

// SaveTask saves a task (or tasks) to the store and returns the first error
// encountered. It does not commit to disk, use `Commit` to do that.
func (j *JSONStore) SaveTask(ts ...*Task) error {
	j.taskLock.Lock()
	defer j.taskLock.Unlock()

	for _, t := range ts {
		// assign the Task an ID if it doesn't have one already
		if t.ID == "" {
			t.ID = goremutake.Encode(j.NextID())
		}

		j.Tasks[t.ID] = t
	}

	return nil
}

// GetTask returns the task for an ID or ErrNoSuchTask
func (j *JSONStore) GetTask(id string) (*Task, error) {
	j.taskLock.RLock()
	defer j.taskLock.RUnlock()

	task, present := j.Tasks[id]
	if !present {
		return nil, ErrNoSuchTask
	}

	return task, nil
}

// DeleteTask deletes a task (or tasks.) It will return an error if any of the
// provided IDs does not exist.
func (j *JSONStore) DeleteTask(ids ...string) error {
	// make sure we have all those tasks before we delete them
	j.taskLock.Lock()
	defer j.taskLock.Unlock()

	for _, id := range ids {
		if _, ok := j.Tasks[id]; !ok {
			return ErrNoSuchTask
		}
	}

	// now we know we have all of them, delete!
	for _, id := range ids {
		delete(j.Tasks, id)
	}

	return nil
}

// FilterTasks taskes a predicate to filter tasks against
func (j *JSONStore) FilterTasks(pred func(*Task) bool) ([]*Task, error) {
	tasks := []*Task{}

	// TODO: this could be parallelized pretty easily
	for _, t := range j.Tasks {
		if pred(t) {
			tasks = append(tasks, t)
		}
	}

	return tasks, nil
}

// SetMeta sets a single meta K/V pair
func (j *JSONStore) SetMeta(k, v string) error {
	j.metaLock.Lock()
	defer j.metaLock.Unlock()

	j.Meta[k] = v
	return nil
}

// GetMeta gets the value of a pair or returns ErrNoSuchKey
func (j *JSONStore) GetMeta(k string) (string, error) {
	j.metaLock.RLock()
	defer j.metaLock.RUnlock()

	v, ok := j.Meta[k]
	if !ok {
		return v, ErrNoSuchKey
	}

	return v, nil
}
