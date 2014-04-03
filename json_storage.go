package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

var ErrNoSuch = errors.New("no such file or directory")

type JSONStore struct {
	filename string           `json:"-"`
	CurID    int              `json:"cur_id"`
	Tasks    map[string]*Task `json:"tasks"`
}

func NewJSONStore(filename string) (*JSONStore, error) {
	store := &JSONStore{filename: filename}

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

func (j *JSONStore) NextID() int {
	j.CurID += 1
	return j.CurID
}

func (j *JSONStore) Commit() error {
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
