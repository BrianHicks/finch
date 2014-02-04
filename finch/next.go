package main

import (
	commander "code.google.com/p/go-commander"
	"fmt"
	"github.com/BrianHicks/finch/core"

	"log"
)

func GetNext(tdb *core.TaskStore, args []string) (*core.Task, error) {
	return tdb.GetNextSelected()
}

var Next *commander.Command = &commander.Command{
	UsageLine: "next",
	Short:     "display the current task",
	Long:      `show the currently active task. This will be the most recently added selected task. To select tasks, run "select".`,
	Run: func(cmd *commander.Command, args []string) {
		tdb, err := getTaskStore()
		defer tdb.Close()
		if err != nil {
			log.Fatalf("Error opening Task database: %s\n", err)
		}

		task, err := GetNext(tdb, args)
		if err != nil {
			log.Fatalf("Error getting task: %s\n", err)
		}

		fmt.Printf("%s\n", task.Description)
	},
}
