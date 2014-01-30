package main

import (
	commander "code.google.com/p/go-commander"
	"fmt"
	"github.com/BrianHicks/finch"
	"strings"
	"time"

	"log"
)

func Adder(tdb *finch.TaskDB, args []string) (*finch.Task, error) {
	task := finch.NewTask(strings.Join(args, " "), time.Now())
	err := tdb.PutTasks(task)
	if err != nil {
		return task, err
	}

	return task, nil
}

var Add *commander.Command = &commander.Command{
	UsageLine: "add",
	Short:     "add a new task to the task database",
	Long:      "add a new task to the task database",
	Run: func(cmd *commander.Command, args []string) {
		tdb, err := getTaskDB()
		defer tdb.Close()
		if err != nil {
			log.Fatalf("Error opening Task database: %s\n", err)
		}

		task, err := Adder(tdb, args)
		if err != nil {
			log.Fatalf("Error adding task: %s\n", err)
		}

		fmt.Printf("Added \"%s\"\n", task.Description)
	},
}
