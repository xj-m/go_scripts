// Package schedule ...
package schedule

import (
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xj-m/go_scripts/file"
)

// func get_times_from_daily_folder
func getFp2timeUnderDir(scheduleDirName string) map[string]time.Time {
	todoFiles, err := file.GetAllFilesWithExtension(scheduleDirName, []string{".todo"})
	if err != nil {
		panic(err)
	}

	// for each file, if contains time string "%Y-%m-%d", add to times
	fp2times := map[string]time.Time{}
	for _, todoFile := range todoFiles {
		t, err := file.ExtractTimeFromFileName(todoFile)
		if err != nil {
			logrus.Error(err)
			continue
		}
		fp2times[todoFile] = t
	}

	return fp2times
}

func ArchiveScheduleTodoFiles(scheduleDirName, archiveDirName string) error {
	for fp, t := range getFp2timeUnderDir(scheduleDirName) {
		if t.Before(file.TruncateToDay(time.Now())) {
			// archive
			year, mon, _ := file.GetYearMonDay(t)
			// create folder
			archivedFolder := filepath.Join(archiveDirName, year, mon)
			err := file.MkdirIfNotExist(archivedFolder)
			if err != nil {
				return err
			}
			archiveFilePath := filepath.Join(archivedFolder, filepath.Base(fp))
			// move file
			err = file.MoveFile(fp, archiveFilePath)
			if err != nil {
				return err
			}
		}
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

func GetTodayTodoFilePath() string {
	scheduleDirName := "schedule"
	return filepath.Join(scheduleDirName, TimeToFileName(time.Now()))
}


func GetTodaySymlink() string {
	scheduleDirName := "schedule"
	return filepath.Join(scheduleDirName, TimeToFileName(time.Now()))
}

