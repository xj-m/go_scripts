// package main
package main

import (
	"path/filepath"
	"time"

	"github.com/xj-m/go_scripts/file"
	"github.com/xj-m/go_scripts/log"
	"github.com/xj-m/go_scripts/task/schedule"
)

func main() {
	err := schedule.ArchiveScheduleTodoFiles(schedule.ScheduleDirPath, schedule.ArchivedDirName)
	if err != nil {
		panic(err)
	}

	// * create today+tmr schedule
	todayFilePath, _ := schedule.GetTodayTodoFilePath()
	for _, newFilePath := range []string{
		todayFilePath,
	} {
		if !file.IsFileExist(newFilePath) {
			// * if today is sunday, use sunday template
			templatePath := schedule.TemplateFilePath
			switch time.Now().Weekday() {
			case time.Sunday:
				templatePath = schedule.SundayTemplateFilePath
				if time.Now().AddDate(0, 0, 7).Month() != time.Now().Month() {
					templatePath = schedule.LastSundayInMonthTemplateFilePath
				}
			}

			// * copy template file to today.todo and tmr.todo
			log.GetLogger(nil).Infof("[task/daily] start to copy template file %s to %s", templatePath, newFilePath)
			if copyErr := file.CopyFileWithIO(templatePath, newFilePath); copyErr != nil {
				panic(copyErr)
			}
		}
	}

	// * create "today.todo" as symlink to today's schedule
	todaySymlink := schedule.GetTodaySymlink()
	log.GetLogger(nil).Infof("[task/daily] created symlink (%v) to (%v)", todaySymlink, todayFilePath)
	if !file.IsFileExist(todaySymlink) {
		if symlinkErr := file.CreateSymlink(filepath.Base(todayFilePath), todaySymlink); symlinkErr != nil {
			panic(symlinkErr)
		}
	}

	// * sort and overwrite todo.todo
	log.GetLogger(nil).Infof("[task/daily] start to sort and overwrite (%v)", schedule.MainTodoFilePath)
	err = schedule.SortAndOverWriteTaskFile(schedule.MainTodoFilePath)
	if err != nil {
		panic(err)
	}

	// * cmdline use "code" open today.todo schedule
	log.GetLogger(nil).Infof("[task/daily] start to open (%v)", todayFilePath)
	err = file.RunCmd("code", todayFilePath)
	if err != nil {
		panic(err)
	}
}
