// package main
package main

import (
	"github.com/xj-m/go_scripts/file"
	"github.com/xj-m/go_scripts/log"
	"github.com/xj-m/go_scripts/photo"
)

func main() {
	todoDir := "todo_gif"
	tmpTodoDir := "tmp_todo_gif"
	dstDir := "generated_gif"
	gifWork(todoDir, tmpTodoDir, dstDir)

	// // get all sub-dir under todoDir
	// subDirs, err := file.GetAllSubDirs(todoDir)
	// if err != nil {
	// 	panic(err)
	// }
	// logrus.Info("subDirs: ", subDirs)

	// for _, subDir := range subDirs {
	// 	gifWork(
	// 		subDir,
	// 		tmpTodoDir,
	// 		dstDir,
	// 	)
	// }
}

func gifWork(todoDir, tmpTodoDir, dstDir string) {
	file.MkdirIfNotExist(dstDir)
	photo.BulkCompressPhoto(todoDir, tmpTodoDir)

	tmpFiles, err := file.GetAllFilesWithExtension(tmpTodoDir, []string{".jpg", ".jpeg"})
	if err != nil {
		panic(err)
	}

	fp, err := photo.MakeGif(tmpFiles, dstDir)
	if err != nil {
		panic(err)
	}

	originalFiles, err := file.GetAllFilesWithExtension(todoDir, []string{".jpg", ".jpeg"})
	if err != nil {
		panic(err)
	}
	// for outputGifFileName, change the create time and modify time to the first file's create time and modify time
	err = file.MatchFileTime(originalFiles[0], fp)

	log.GetLogger(nil).Info("gif generated at: ", fp)

	// del dir tmpTodoDir
	err = file.DeleteDir(tmpTodoDir)
	if err != nil {
		panic(err)
	}
}
