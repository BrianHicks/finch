package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"sync"

	"github.com/cryptix/goremutake"
)

var ErrNoSuch = errors.New("no such file or directory")

type JSONStore struct {
	filename string

	CurID uint `json:"cur_id"`

	Tasks    map[string]*Task `json:"tasks"`
	taskLock *sync.RWMutex

	Meta     map[string]string `json:"meta"`
	metaLock *sync.RWMutex
}

func NewJSONStore(filename string) (*JSONStore, error) {
	store := &JSONStore{
		filename: filename,
		taskLock: new(sync.RWMutex),
		metaLock: new(sync.RWMutex),
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		if perr, ok := err.(*os.PathError); ok {
			if perr.Err.Error() == ErrNoSuch.Error() {
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

func (j *JSONStore) NextID() uint {
	j.CurID += 1
	return j.CurID
}

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

func (j *JSONStore) GetTask(id string) (*Task, error) {
	j.taskLock.RLock()
	defer j.taskLock.RUnlock()

	task, present := j.Tasks[id]
	if !present {
		return nil, NoSuchTask
	}

	return task, nil
}

func (j *JSONStore) DeleteTask(ids ...string) error {
	// make sure we have all those tasks before we delete them
	j.taskLock.Lock()
	defer j.taskLock.Unlock()

	for _, id := range ids {
		if _, ok := j.Tasks[id]; !ok {
			return NoSuchTask
		}
	}

	// now we know we have all of them, delete!
	for _, id := range ids {
		delete(j.Tasks, id)
	}

	return nil
}

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

func (j *JSONStore) SetMeta(k, v string) error {
	j.metaLock.Lock()
	defer j.metaLock.Unlock()

	j.Meta[k] = v
	return nil
}

func (j *JSONStore) GetMeta(k string) (string, error) {
	j.metaLock.RLock()
	defer j.metaLock.RUnlock()

	v, ok := j.Meta[k]
	if !ok {
		return v, NoSuchKey
	}

	return v, nil
}
