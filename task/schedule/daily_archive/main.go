// package main
package main

import (
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/xj-m/go_scripts/file"
	"github.com/xj-m/go_scripts/task/schedule"
)

func main() {
	err := schedule.ArchiveScheduleTodoFiles(schedule.ScheduleDirPath, schedule.ArchivedDirName)
	if err != nil {
		panic(err)
	}

	// create today+tmr schedule
	todayFilePath, _ := schedule.GetTodayTodoFilePath()
	tmrFilename, _ := schedule.GetTmrTodoFilePath()
	for _, newFilePath := range []string{
		todayFilePath,
		tmrFilename,
	} {
		if !file.IsFileExist(newFilePath) {
			logrus.Infof("[task/daily] start to create new schedule file %s", newFilePath)
			if copyErr := file.CopyFileWithIO(schedule.TemplateFilePath, newFilePath); copyErr != nil {
				panic(copyErr)
			}
		}
	}

	// create "today.todo" as symlink to today's schedule
	todaySymlink := schedule.GetTodaySymlink()
	logrus.Infof("[task/daily] created symlink (%v) to (%v)", todaySymlink, todayFilePath)
	if !file.IsFileExist(todaySymlink) {
		if symlinkErr := file.CreateSymlink(filepath.Base(todayFilePath), todaySymlink); symlinkErr != nil {
			panic(symlinkErr)
		}
	}

	// sort and overwrite todo.todo
	logrus.Infof("[task/daily] start to sort and overwrite (%v)", schedule.MainTodoFilePath)
	err = schedule.SortAndOverWriteTaskFile(schedule.MainTodoFilePath)
	if err != nil {
		panic(err)
	}

	// cmdline use "code" open today.todo schedule
	logrus.Infof("[task/daily] start to open (%v)", todayFilePath)
	err = file.RunCmd("code", todayFilePath)
	if err != nil {
		panic(err)
	}
}
