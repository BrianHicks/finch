package main

import (
	commander "code.google.com/p/go-commander"
	"fmt"

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
		Version,
	},
}

func main() {
	err := commands.Run(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
