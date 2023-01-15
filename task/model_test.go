package task

import (
	"io/ioutil"
	"testing"
)

func TestTask_ToContent(t *testing.T) {
	taskA, err := ParseTaskFromTodoFile("2023-01-12.todo")
	if err != nil {
		t.Fatal(err)
	}
	taskB, err := ParseTaskFromTodoFile("2023-01-13.todo")
	if err != nil {
		t.Fatal(err)
	}
	taskA = *taskA.TaskName2task["all"]
	taskB = *taskB.TaskName2task["all"]
	got := MergeTasks(taskA, taskB)
	WriteStringToFile(got.ToContent(), "ut_output.todo")
	// write string `got` to "ut_output.todo"
}

func WriteStringToFile(s, filename string) error {
	return ioutil.WriteFile(filename, []byte(s), 0o644)
}
