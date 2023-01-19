// package task
package task

func FilterHighPriority(item Item) bool {
	return hasKey(item.TagK2v, "high")
}

func FilterNotHighPriority(item Item) bool {
	return !hasKey(item.TagK2v, "high")
}

func FilterNotRoutineArchive(task Task) bool {
	return !contains([]string{"Archive", "routine"}, task.Name)
}

func FilterNotEmptyTask(task Task) bool {
	return len(task.Items) > 0 || len(task.TaskName2task) > 0
}

func FilterNotEmptyItem(item Item) bool {
	return item.Content != ""
}
