// package task
package task

func FilterItemHighPriority(item Item) bool {
	hasParentTaskHighPriority := false
	task := item.ParentTask
	for task != nil {
		if contains(task.TagNames, "high") {
			hasParentTaskHighPriority = true
			break
		}
		task = task.Parent
	}
	return hasKey(item.TagK2v, "high") || hasParentTaskHighPriority
}

func FilterItemNotHighPriority(item Item) bool {
	return !FilterItemHighPriority(item)
}

func FilterItemNotEmpty(item Item) bool {
	return item.Content != ""
}

func FilterTaskRoutineArchive(task Task) bool {
	hasParentRoutineArchive := false
	currentTask := &task
	for currentTask != nil {
		if contains([]string{"Archive", "routine"}, currentTask.Name) {
			hasParentRoutineArchive = true
			break
		}
		currentTask = currentTask.Parent
	}
	return hasParentRoutineArchive
}

func FilterTaskNotRoutineArchive(task Task) bool {
	return !FilterTaskRoutineArchive(task)
}

func FilterTaskNotEmpty(task Task) bool {
	return len(task.Items) > 0 || len(task.TaskName2task) > 0
}
