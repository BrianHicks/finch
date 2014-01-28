package finch

func reverseTaskSlice(orig []*Task) []*Task {
	tasks := []*Task{}

	for i := len(orig) - 1; i >= 0; i-- {
		tasks = append(tasks, orig[i])
	}

	return tasks
}
