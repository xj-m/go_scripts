// package main
package main

import (
	"os"
	"time"

	"github.com/xj-m/go_scripts/file"
	"github.com/xj-m/go_scripts/task"
)

var (
	SCHEDULE_DIR_NAME = "schedule"
	ARCHIVED_DIR_NAME = "archived/todo"
	TMP_FILE_NAME     = "temp-daily.todo"
	MAIN_TODO_FILE    = "todo.todo"
)

func main() {
	curTodoFilePath := SCHEDULE_DIR_NAME + "/" + file.TruncateToDay(time.Now()).Format("2006-01-02") + ".todo"
	// if curTodoFilePath not exist, panic
	if _, err := os.Stat(curTodoFilePath); err != nil {
		panic(err)
	}

	srcFilePath := curTodoFilePath
	dstFilePath := MAIN_TODO_FILE

	srcTask, err := task.ParseTaskFromTodoFile(srcFilePath)
	if err != nil {
		panic(err)
	}
	// remove parts that doesn't want to be merged
	srcTask.TaskName2task["all"].RemoveTaskByNames([]string{
		"routine",
	})
	srcTask.RemoveTaskByNames([]string{
		"Archive",
	})

	// read dst
	dstTask, err := task.ParseTaskFromTodoFile(dstFilePath)
	if err != nil {
		panic(err)
	}
	// merge
	mergedTask := task.MergeTasks(srcTask, dstTask)
	// write to dst
	err = file.OverWriteFile(dstFilePath, mergedTask.ToContent())
	if err != nil {
		panic(err)
	}
}
