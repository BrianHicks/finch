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
}

func NewJSONStore(filename string) (*JSONStore, error) {
	store := &JSONStore{
		filename: filename,
		taskLock: new(sync.RWMutex),
		Tasks:    map[string]*Task{},
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		if perr, ok := err.(*os.PathError); ok {
			if perr.Err.Error() == ErrNoSuch.Error() {
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

func (j *JSONStore) SaveTask(t *Task) error {
	// assign the Task an ID if it doesn't have one already
	if t.ID == "" {
		t.ID = goremutake.Encode(j.NextID())
	}

	j.taskLock.Lock()
	defer j.taskLock.Unlock()

	j.Tasks[t.ID] = t

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
