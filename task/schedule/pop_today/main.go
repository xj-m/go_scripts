// package main
package main

import (
	"github.com/xj-m/go_scripts/file"
	"github.com/xj-m/go_scripts/task"
	"github.com/xj-m/go_scripts/task/schedule"
)

func main() {
	todayFilePath, err := schedule.GetTodayTodoFilePath()
	if err != nil {
		panic(err)
	}

	// read src
	srcFilePath := schedule.MainTodoFilePath
	srcTask, err := task.ParseTaskFromTodoFile(srcFilePath)
	if err != nil {
		panic(err)
	}

	// remove parts that doesn't want to be merged
	srcTask.TaskName2task["all"].FilterItems(
		task.FilterItemHighPriority,
		task.FilterItemNotEmpty,
	)
	srcTask.Filter(
		task.FilterTaskNotRoutineArchive,
		task.FilterTaskNotEmpty,
	)

	// read dst
	dstTask, err := task.ParseTaskFromTodoFile(todayFilePath)
	if err != nil {
		panic(err)
	}

	// merge
	mergedTask := task.MergeTasks(srcTask, dstTask)

	// write to dst
	err = file.OverWriteFile(todayFilePath, mergedTask.ToContent())
	if err != nil {
		panic(err)
	}

	// update src with high priority task removed
	srcTask, _ = task.ParseTaskFromTodoFile(srcFilePath)
	srcTask.TaskName2task["all"].FilterItems(task.FilterItemNotHighPriority, task.FilterItemNotEmpty)
	srcTask.Filter(
		task.FilterTaskNotEmpty,
	)
	err = file.OverWriteFile(srcFilePath, srcTask.ToContent())
	if err != nil {
		panic(err)
	}
}
