// package task
package task

func HighPriorityFilter(item Item) bool {
	return hasKey(item.TagK2v, "high")
}

func FilterRoutineArchive(task Task) bool {
	return !contains([]string{"Archive", "routine"}, task.Name)
}
