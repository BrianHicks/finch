package main

import (
	commander "code.google.com/p/go-commander"
	"github.com/BrianHicks/finch"
	"time"

	"log"
)

func Delayer(tdb *finch.TaskDB, args []string) (*finch.Task, error) {
	task, err := tdb.GetNextSelected()
	if err != nil {
		return task, err
	}
	oldKey := task.Key()

	task.Timestamp = time.Now()

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
		tdb, err := getTaskDB()
		defer tdb.Close()
		if err != nil {
			log.Fatalf("Error opening Task database: %s\n", err)
		}

		_, err = Delayer(tdb, args)
		if err != nil {
			log.Fatalf("Error delaying task: %s\n", err)
		}
	},
}
