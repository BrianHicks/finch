package main

import (
	commander "code.google.com/p/go-commander"

	"log"
)

var Done *commander.Command = &commander.Command{
	UsageLine: "done",
	Short:     "mark current task as done",
	Long:      `mark the current task (from "next") as done. If you're not actually *done* with this task, use "delay"`,
	Run: func(cmd *commander.Command, args []string) {
		log.Printf("done: %+v\n", args)
	},
}
