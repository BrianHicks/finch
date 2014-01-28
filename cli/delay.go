package main

import (
	commander "code.google.com/p/go-commander"

	"log"
)

var Delay *commander.Command = &commander.Command{
	UsageLine: "delay",
	Short:     "mark current task as delayed",
	Long: `mark the current task (from "next") as delayed.

This will re-enter this task at the end of the database.`,
	Run: func(cmd *commander.Command, args []string) {
		log.Printf("delayed: %+v\n", args)
	},
}
