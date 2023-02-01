package task

import (
	"io/ioutil"
	"reflect"
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

func TestTask_SubTask(t *testing.T) {
	type fields struct {
		TaskName2task map[string]*Task
		TaskNames     []string
		Items         Items
		Level         int
		Name          string
		Parent        *Task
		TagK2v        tagK2v
		TagNames      []string
		Status        string
		Comments      []string
	}
	type args struct {
		taskToRemove *Task
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Task
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Task{
				TaskName2task: tt.fields.TaskName2task,
				TaskNames:     tt.fields.TaskNames,
				Items:         tt.fields.Items,
				Level:         tt.fields.Level,
				Name:          tt.fields.Name,
				Parent:        tt.fields.Parent,
				TagK2v:        tt.fields.TagK2v,
				TagNames:      tt.fields.TagNames,
				Status:        tt.fields.Status,
				Comments:      tt.fields.Comments,
			}
			if got := tr.SubTask(tt.args.taskToRemove); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Task.SubTask() = %v, want %v", got, tt.want)
			}
		})
	}
}
