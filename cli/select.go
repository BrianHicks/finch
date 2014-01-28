package main

import (
	commander "code.google.com/p/go-commander"

	"log"
)

var Select *commander.Command = &commander.Command{
	UsageLine: "select",
	Short:     "select tasks to work on",
	Long: `select tasks to work on, as a part of a new chain.

Selected tasks will be marked as such, then you should complete them in the
*reverse* order you selected them. You can use the "next" command for this.`,
	Run: func(cmd *commander.Command, args []string) {
		log.Printf("selected: %+v\n", args)
	},
}
