package main

import (
	commander "code.google.com/p/go-commander"
	"fmt"
	"github.com/BrianHicks/finch/core"
	"os"
	"strconv"
	"text/template"

	"log"
)

var (
	taskTmpl = template.Must(template.New("task").Parse("{{.Count}}: {{.Task.Description}} {{if .Selected}}(selected){{end}}\n"))
)

type IterTask struct {
	Count    int
	Task     *core.Task
	Selected bool
}

func Selector(tdb *core.TaskDB, args []string) ([]*core.Task, error) {
	tasks, err := tdb.GetPendingTasks()
	if err != nil {
		return tasks, err
	}

	amtTasks := len(tasks)

	if len(args) == 0 {
		// print tasks so the user can see and select them
		for i := 0; i < amtTasks; i++ {
			task := tasks[i]
			selected, ok := task.Attrs[core.TagSelected]

			err := taskTmpl.Execute(os.Stdout, IterTask{i, task, selected && ok})
			if err != nil {
				break
				fmt.Printf("Error: %s", err.Error())
			}
		}

		return []*core.Task{}, nil
	} else {
		// select tasks and print status for each one
		selected := []*core.Task{}

		for i := 0; i < len(args); i++ {
			fmt.Printf("selecting \"%s\"... ", args[i])

			taskNum, err := strconv.Atoi(args[i])
			if err != nil {
				fmt.Print("was not digits!\n")
				continue
			}

			if taskNum > amtTasks-1 {
				fmt.Printf("there is no task %d\n", amtTasks)
				continue
			}

			task := tasks[taskNum]
			task.Attrs[core.TagSelected] = true
			selected = append(selected, task)
			fmt.Printf("selected \"%s\"\n", task.Description)
		}

		if len(selected) > 0 {
			err := tdb.PutTasks(selected...)
			if err != nil {
				return selected, err
			} else {
				fmt.Printf("Wrote %d tasks to DB\n", len(selected))
			}
		}

		return selected, nil
	}

	return []*core.Task{}, nil
}

var Select *commander.Command = &commander.Command{
	UsageLine: "select",
	Short:     "select tasks to work on",
	Long: `select tasks to work on, as a part of a new chain.

Selected tasks will be marked as such, then you should complete them in the
*reverse* order you selected them. You can use the "next" command for this.

To select tasks, call "select" once to view a list of tasks. Then you can run
"select" with the numbers of the tasks you want to select them.`,
	Run: func(cmd *commander.Command, args []string) {
		tdb, err := getTaskDB()
		defer tdb.Close()
		if err != nil {
			log.Fatalf("Error opening Task database: %s\n", err)
		}

		_, err = Selector(tdb, args)
		if err != nil {
			log.Fatalf("Error selecting tasks: %s\n", err)
		}
	},
}
