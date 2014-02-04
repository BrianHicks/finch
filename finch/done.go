package main

import (
	commander "code.google.com/p/go-commander"
	"fmt"
	"github.com/BrianHicks/finch/core"

	"log"
)

func MarkDone(tdb *core.TaskStore, args []string) (*core.Task, error) {
	task, err := tdb.GetNextSelected()
	if err != nil {
		return task, err
	}

	task.Attrs[core.TagSelected] = false
	task.Attrs[core.TagPending] = false

	// TODO: callbacks to other programs here?

	err = tdb.PutTasks(task)
	if err != nil {
		return task, err
	}

	return task, nil
}

var Done *commander.Command = &commander.Command{
	UsageLine: "done",
	Short:     "mark current task as done",
	Long:      `mark the current task (from "next") as done. If you're not actually *done* with this task, use "delay"`,
	Run: func(cmd *commander.Command, args []string) {
		tdb, err := getTaskStore()
		defer tdb.Store.Close()
		if err != nil {
			log.Fatalf("Error opening Task database: %s\n", err)
		}

		task, err := MarkDone(tdb, args)
		if err != nil {
			log.Fatalf("Error marking task done: %s\n", err)
		}

		fmt.Printf("Marked \"%s\" done\n", task.Description)
	},
}
