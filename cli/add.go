package main

import (
	commander "code.google.com/p/go-commander"

	"log"
)

var Add *commander.Command = &commander.Command{
	UsageLine: "add",
	Short:     "add a new task to the task database",
	Long:      "add a new task to the task database",
	Run: func(cmd *commander.Command, args []string) {
		log.Printf("added: %+v\n", args)
	},
}
