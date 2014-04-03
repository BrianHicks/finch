package main

import (
	"encoding/json"
	"io/ioutil"
)

type JSONStore struct {
	filename string
	CurID    int
	Tasks    []*Task
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
