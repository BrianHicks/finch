package main

import (
	"bytes"
	"text/template"
	"time"
)

var taskTmpl = template.Must(template.New("task").Parse(`{{.ID}}: {{.Desc}}{{if .Selected}} (*){{end}}{{if .Done}} (done){{end}}`))

type Task struct {
	ID       string
	Desc     string
	Active   time.Time
	Done     bool
	Selected bool
	Repeat   time.Duration
}

func (t *Task) String() string {
	var s bytes.Buffer

	err := taskTmpl.Execute(&s, t)
	if err != nil {
		panic(err)
	}

	return s.String()
}

func (t *Task) MarkDone() {
	t.Selected = false
	if t.Repeat == 0 {
		t.Done = true
	} else {
		t.Active = time.Now().Add(t.Repeat)
	}
}

type ByActive []*Task

func (a ByActive) Len() int           { return len(a) }
func (a ByActive) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByActive) Less(i, j int) bool { return a[i].Active.Before(a[j].Active) }
