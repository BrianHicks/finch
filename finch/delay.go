package main

import (
	commander "code.google.com/p/go-commander"
	"fmt"
	"github.com/BrianHicks/finch/core"
	"time"

	"log"
)

func Delayer(tdb *core.TaskStore, args []string) (*core.Task, error) {
	task, err := tdb.GetNextSelected()
	if err != nil {
		return task, err
	}
	oldKey := task.Key()

	task.Timestamp = time.Now()
	task.Attrs[core.TagSelected] = false

	err = tdb.MoveTask(oldKey, task)
	if err != nil {
		return task, err
	}

	return task, nil
}

var Delay *commander.Command = &commander.Command{
	UsageLine: "delay",
	Short:     "mark current task as delayed",
	Long: `mark the current task (from "next") as delayed.

This will re-enter this task at the end of the database.`,
	Run: func(cmd *commander.Command, args []string) {
		tdb, err := getTaskStore()
		defer tdb.Store.Close()
		if err != nil {
			log.Fatalf("Error opening Task database: %s\n", err)
		}

		task, err := Delayer(tdb, args)
		if err != nil {
			log.Fatalf("Error delaying task: %s\n", err)
		}

		fmt.Printf("Delayed \"%s\"\n", task.Description)
	},
}
