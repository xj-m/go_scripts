// package main
package main

import (
	"fmt"
	"path/filepath"

	"github.com/xj-m/go_scripts/file"
	"github.com/xj-m/go_scripts/task/schedule"
)

var (
	scheduleDirName = "schedule"
	archivedDirName = "archived/todo"
	tmpFilePath     = fmt.Sprintf("%s/%s", scheduleDirName, "temp-daily.todo")
)

func main() {
	err := schedule.ArchiveScheduleTodoFiles(scheduleDirName, archivedDirName)
	if err != nil {
		panic(err)
	}

	// create today+tmr schedule
	todayFilePath := schedule.GetTodayTodoFilePath()
	tmrFilename := schedule.GetTmrTodoFilePath()
	for _, newFilePath := range []string{
		todayFilePath,
		tmrFilename,
	} {
		if !file.IsFileExist(newFilePath) {
			if copyErr := file.CopyFileWithIO(tmpFilePath, newFilePath); copyErr != nil {
				panic(copyErr)
			}
		}
	}

	// create "today.todo" as symlink to today's schedule
	todaySymlink := schedule.GetTodaySymlink()
	if !file.IsFileExist(todaySymlink) {
		if symlinkErr := file.CreateSymlink(filepath.Base(todayFilePath), todaySymlink); symlinkErr != nil {
			panic(symlinkErr)
		}
	}

	// sort and overwrite todo.todo
	err = schedule.SortAndOverWriteTaskFile("todo.todo")
	if err != nil {
		panic(err)
	}

	// cmdline use "code" open today.todo schedule
	err = file.RunCmd("code", todayFilePath)
	if err != nil {
		panic(err)
	}
}
