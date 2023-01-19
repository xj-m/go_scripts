// package main
package main

import (
	"fmt"
	"time"

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
	tmrFilename := addScheduleFolder(schedule.TimeToFileName(time.Now().Add(24 * time.Hour)))
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
	todaySymlink := addScheduleFolder("today.todo")
	if !file.IsFileExist(todaySymlink) {
		if symlinkErr := file.CreateSymlink(todayFileName, todaySymlink); symlinkErr != nil {
			panic(symlinkErr)
		}
	}
	// cmdline use "code" open today.todo schedule
	err = file.RunCmd("code", todayFilePath)
	if err != nil {
		panic(err)
	}
}
