// Package schedule ...
package schedule

import "testing"

func TestMoveTaskAndOverwriteBoth(t *testing.T) {
	srcFilePath := "2023-01-18.todo"
	dstFilePath := "todo.todo"
	MoveTaskAndOverwriteBoth(srcFilePath, dstFilePath)
}
