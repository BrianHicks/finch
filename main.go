package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/BrianHicks/finch/duration"
	"github.com/codegangsta/cli"
)

var (
	commands = []cli.Command{
		{
			Name:      "task-add",
			ShortName: "add",
			Usage:     "add a task to the list (from concatenation of args)",
			Action: func(c *cli.Context) {
				coord, err := taskCoordinatorFromContext(c)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error opening tasks: %s\n", err.Error())
					return
				}

				if !c.Args().Present() {
					fmt.Fprint(os.Stderr, "add requires a description\n")
					return
				}

				// TODO: parse for repetition/delay?
				task := coord.Add(strings.Join(c.Args(), " "))
				err = coord.Save(task)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error saving task: %s\n", err.Error())
					return
				}

				err = coord.Close()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error saving tasks: %s\n", err.Error())
				}

				fmt.Fprintf(os.Stdout, "Added: %s\n", task)
			},
		},
		{
			Name:      "task-available",
			ShortName: "available",
			Usage:     "view available tasks",
			Action: func(c *cli.Context) {
				coord, err := taskCoordinatorFromContext(c)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error opening tasks: %s\n", err.Error())
					return
				}

				tasks, err := coord.Available()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error reading tasks: %s\n", err.Error())
					return
				}

				for _, task := range tasks {
					fmt.Println(task)
				}
			},
		},
		{
			Name:      "task-selected",
			ShortName: "current",
			Usage:     "view selected tasks",
			Action: func(c *cli.Context) {
				coord, err := taskCoordinatorFromContext(c)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error opening tasks: %s\n", err.Error())
					return
				}

				tasks, err := coord.Selected()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error reading tasks: %s\n", err.Error())
					return
				}

				for _, task := range tasks {
					fmt.Println(task)
				}
			},
		},
		{
			Name:      "task-next",
			ShortName: "next",
			Usage:     "see the next selected task (selected, but only 1)",
			Action: func(c *cli.Context) {
				coord, err := taskCoordinatorFromContext(c)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error opening tasks: %s\n", err.Error())
					return
				}

				task, err := coord.NextSelected()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error reading tasks: %s\n", err.Error())
					return
				}

				fmt.Println(task)
			},
		},
		{
			Name:      "task-select",
			ShortName: "select",
			Usage:     "select tasks",
			Action: func(c *cli.Context) {
				coord, err := taskCoordinatorFromContext(c)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error opening tasks: %s\n", err.Error())
					return
				}

				err = coord.Select(c.Args()...)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error selecting tasks: %s\n", err.Error())
				}

				err = coord.Close()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error saving selection: %s\n", err.Error())
				}

				fmt.Fprintf(os.Stdout, "Saved: %s\n", strings.Join(c.Args(), ", "))
			},
		},
		{
			Name:      "task-delay",
			ShortName: "delay",
			Usage:     "delay a task",
			Flags: []cli.Flag{
				cli.StringFlag{"until, u", "", "delay until this date (ISO-8601 specified)"},
			},
			Action: func(c *cli.Context) {
				coord, err := taskCoordinatorFromContext(c)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error opening tasks: %s\n", err.Error())
					return
				}

				// get specified task or current task
				var task *Task
				if c.Args().Present() {
					task, err = coord.Get(c.Args().First())
				} else {
					task, err = coord.NextSelected()
				}
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error getting next task: %s\n", err.Error())
					return
				}

				fmt.Fprintf(os.Stdout, "Delaying: %s\n", task)

				// parse until
				until := time.Now()
				if c.String("until") != "" {
					dur, err := duration.FromString(c.String("until"))
					if err != nil {
						fmt.Fprintf(os.Stderr, "Error parsing until: %s\n", err.Error())
						return
					}

					until = until.Add(dur.ToDuration())
				}

				err = coord.Delay(task.ID, until)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error delaying %s: %s\n", task.ID, err.Error())
					return
				}

				err = coord.Close()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error saving %s: %s\n", task.ID, err.Error())
					return
				}

				fmt.Fprintf(os.Stdout, "Delayed until %s\n", until)
			},
		},
		{
			Name:      "task-done",
			ShortName: "done",
			Usage:     "mark a task done",
			Action: func(c *cli.Context) {
				coord, err := taskCoordinatorFromContext(c)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error opening tasks: %s\n", err.Error())
					return
				}

				// get specified task or current task
				var task *Task
				if c.Args().Present() {
					task, err = coord.Get(c.Args().First())
				} else {
					task, err = coord.NextSelected()
				}
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error getting next task: %s\n", err.Error())
					return
				}

				err = coord.MarkDone(task.ID)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error marking task done: %s\n", err.Error())
					return
				}

				err = coord.Close()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error writing tasks: %s\n", err.Error())
					return
				}

				fmt.Fprintf(os.Stdout, "Marked done: %s\n", task)
			},
		},
		{
			Name:      "task-delete",
			ShortName: "delete",
			Usage:     "delete a task",
			Action: func(c *cli.Context) {
				if !c.Args().Present() {
					fmt.Fprintln(os.Stderr, "Error: Need a task to delete.")
					return
				}

				coord, err := taskCoordinatorFromContext(c)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error opening tasks: %s\n", err.Error())
					return
				}

				task, err := coord.Get(c.Args().First())
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error getting next task: %s\n", err.Error())
					return
				}

				err = coord.Delete(task.ID)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error deleting done: %s\n", err.Error())
					return
				}

				err = coord.Close()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error writing tasks: %s\n", err.Error())
					return
				}

				fmt.Fprintf(os.Stdout, "Deleted: %s\n", task)
			},
		},
	}
	flags = []cli.Flag{
		cli.StringFlag{"location, l", "~/.finch", "location to store data"},
	}
)

func main() {
	app := cli.NewApp()
	app.Usage = "command line task management"
	app.Version = "0.0.1"
	app.Commands = commands
	app.Flags = flags

	app.Run(os.Args)
}
