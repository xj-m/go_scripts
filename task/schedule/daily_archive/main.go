// package main
package main

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xj-m/go_scripts/file"
	"github.com/xj-m/go_scripts/task/schedule"
)

var (
	SCHEDULE_DIR_NAME = "schedule"
	ARCHIVED_DIR_NAME = "archived/todo"
	TMP_FILE_NAME     = "temp-daily.todo"
)

func main() {
	err := schedule.ArchiveScheduleTodoFiles(SCHEDULE_DIR_NAME, ARCHIVED_DIR_NAME)
	if err != nil {
		panic(err)
	}
	// create today+tmr schedule
	for _, t := range []time.Time{
		file.TruncateToDay(time.Now()),
		file.TruncateToDay(time.Now().Add(24 * time.Hour)),
	} {
		err := schedule.DuplicateFromFile(t, SCHEDULE_DIR_NAME, TMP_FILE_NAME)
		if err != nil {
			logrus.Error(err)
		}
	}
}
