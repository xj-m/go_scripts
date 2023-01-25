// package main
package main

import (
	"github.com/xj-m/go_scripts/file"
	"github.com/xj-m/go_scripts/task/schedule"
)

func main() {
	for _, srcFilePath := range getToCleanupFilePaths() {
		schedule.MoveTaskAndOverwriteBoth(srcFilePath, schedule.MainTodoFilePath)
	}
}

func getToCleanupFilePaths() []string {
	return file.Keys(schedule.GetFp2timeUnderDir(schedule.ScheduleDirPath))
}
