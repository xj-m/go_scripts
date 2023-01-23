// Package schedule ...
package schedule

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/xj-m/go_scripts/file"
	"github.com/xj-m/go_scripts/log"
	"github.com/xj-m/go_scripts/task"
)

var (
	ScheduleDirPath                   = "schedule"
	MainTodoFilePath                  = "todo.todo"
	TemplateFilePath                  = ScheduleDirPath + "/temp-daily.todo"
	SundayTemplateFilePath            = ScheduleDirPath + "/temp-sunday.todo"
	LastSundayInMonthTemplateFilePath = ScheduleDirPath + "/temp-last_sunday.todo"
	ArchivedDirName                   = "archived/todo"
)

// func get_times_from_daily_folder
func GetFp2timeUnderDir(scheduleDirName string) map[string]time.Time {
	todoFiles, err := file.GetAllFilesWithExtension(scheduleDirName, []string{".todo"})
	if err != nil {
		panic(err)
	}

	// for each file, if contains time string "%Y-%m-%d", add to times
	fp2times := map[string]time.Time{}
	for _, todoFile := range todoFiles {
		t, err := file.ExtractTimeFromFileName(todoFile)
		if err != nil {
			log.GetLogger(nil).Error(err)
			continue
		}
		fp2times[todoFile] = t
	}

	return fp2times
}

func ArchiveScheduleTodoFiles(scheduleDirName, archiveDirName string) error {
	for fp, t := range GetFp2timeUnderDir(scheduleDirName) {
		if t.Before(file.TruncateToDay(time.Now())) {
			// move unfinished task to main todo file
			MoveTaskAndOverwriteDst(fp, MainTodoFilePath)

			// archive
			year, mon, _ := file.GetYearMonDay(t)

			// create folder
			archivedFolder := filepath.Join(archiveDirName, year, mon)
			err := file.MkdirIfNotExist(archivedFolder)
			if err != nil {
				return fmt.Errorf("failed to create folder %s: %w", archivedFolder, err)
			}
			archiveFilePath := filepath.Join(archivedFolder, filepath.Base(fp))

			// move file
			err = file.MoveFile(fp, archiveFilePath)
			if err != nil {
				return fmt.Errorf("failed to move file %s to %s: %w", fp, archiveFilePath, err)
			}
		}
	}
	return nil
}

func MoveTaskAndOverwriteDst(srcFilePath, dstFilePath string) error {
	log.GetLogger(nil).Infof("[task] move task from %s to %s", srcFilePath, dstFilePath)

	// parse task from src
	srcTask, err := task.ParseTaskFromTodoFile(srcFilePath)
	if err != nil {
		return err
	}
	srcTask.TaskName2task["all"].Filter(task.FilterNotRoutineArchive)

	// read dst
	dstTask, err := task.ParseTaskFromTodoFile(dstFilePath)
	if err != nil {
		panic(err)
	}

	// create new task with merged items
	mergedTask := task.MergeTasks(srcTask, dstTask)
	mergedTask.SortItemsByPriority()

	// write to dst
	err = file.OverWriteFile(dstFilePath, mergedTask.ToContent())
	if err != nil {
		return err
	}
	return nil
}

func DuplicateFromFile(newFileName string, scheduleFileDir string, tempFileName string) error {
	// new file named like "2020-01-01.todo"
	newFilePath := filepath.Join(scheduleFileDir, newFileName)
	tempFilePath := filepath.Join(scheduleFileDir, tempFileName)

	// if file not exist then copy from fp
	if !file.IsFileExist(newFilePath) {
		err := file.CopyFileWithIO(tempFilePath, newFilePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func TimeToFileName(t time.Time) string {
	t = file.TruncateToDay(t)
	return t.Format("2006-01-02") + ".todo"
}

func GetYesterdayTodoFilePath() (string, error) {
	ret := filepath.Join(ScheduleDirPath, TimeToFileName(time.Now().Add(-24*time.Hour)))
	if !file.IsFileExist(ret) {
		return ret, file.ErrFileNotExist
	}
	return ret, nil
}

func GetTodayTodoFilePath() (string, error) {
	ret := filepath.Join(ScheduleDirPath, TimeToFileName(time.Now()))
	if !file.IsFileExist(ret) {
		return ret, file.ErrFileNotExist
	}
	return ret, nil
}

func GetTmrTodoFilePath() (string, error) {
	ret := filepath.Join(ScheduleDirPath, TimeToFileName(time.Now().Add(24*time.Hour)))
	if !file.IsFileExist(ret) {
		return ret, file.ErrFileNotExist
	}
	return ret, nil
}

func GetTodaySymlink() string {
	todayTodoFilePath, _ := GetTodayTodoFilePath()
	// create "today.todo" symlink pointing to today's todo file
	todaySymlink := filepath.Join(ScheduleDirPath, "today.todo")
	file.CreateSymlink(todayTodoFilePath, todaySymlink)
	return todaySymlink
}

func SortAndOverWriteTaskFile(fp string) error {
	task, _ := task.ParseTaskFromTodoFile(fp)
	task.SortItemsByPriority()
	return file.OverWriteFile(fp, task.ToContent())
}
