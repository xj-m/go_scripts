// package main
package main

import (
	"os"

	"github.com/xj-m/go_scripts/file"
	"github.com/xj-m/go_scripts/task"
	"github.com/xj-m/go_scripts/task/schedule"
)

var mainTodoFileName = "todo.todo"

func main() {
	todayFilePath := schedule.GetTodayTodoFilePath()

	// if curTodoFilePath not exist, panic
	if _, err := os.Stat(todayFilePath); err != nil {
		panic(err)
	}

	srcFilePath := mainTodoFileName

	srcTask, err := task.ParseTaskFromTodoFile(srcFilePath)
	if err != nil {
		panic(err)
	}

	// remove parts that doesn't want to be merged
	srcTask.TaskName2task["all"].FilterItems(
		task.FilterHighPriority,
		task.FilterNotEmptyItem,
	)
	srcTask.Filter(
		task.FilterNotRoutineArchive,
		task.FilterNotEmptyTask,
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
	srcTask.TaskName2task["all"].FilterItems(task.FilterNotHighPriority)
	err = file.OverWriteFile(srcFilePath, srcTask.ToContent())
	if err != nil {
		panic(err)
	}
}
