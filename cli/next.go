package main

import (
	commander "code.google.com/p/go-commander"

	"log"
)

var Next *commander.Command = &commander.Command{
	UsageLine: "next",
	Short:     "display the current task",
	Long:      `show the currently active task. This will be the most recently added selected task. To select tasks, run "select".`,
	Run: func(cmd *commander.Command, args []string) {
		log.Printf("next: %+v\n", args)
	},
}
