package main

import (
	commander "code.google.com/p/go-commander"

	"os"
)

var commands commander.Commander = commander.Commander{
	Name: "finch",
	Commands: []*commander.Command{
		Add,
		Select,
		Next,
		Done,
		Delay,
	},
}

func main() {
	commands.Run(os.Args[1:])
}
