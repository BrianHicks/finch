package main

import (
	commander "code.google.com/p/go-commander"

	"fmt"
)

var Version *commander.Command = &commander.Command{
	UsageLine: "version",
	Short:     "show version number and exit",
	Long:      "show version number and exit",
	Run: func(cmd *commander.Command, args []string) {
		fmt.Println("0.1.1")
	},
}
