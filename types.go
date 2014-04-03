package main

import "time"

type Task struct {
	ID       int
	Desc     string
	Added    time.Time
	Delay    time.Time
	Done     bool
	Selected bool
}
