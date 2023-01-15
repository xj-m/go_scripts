package task

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func TestMergeTasks(t *testing.T) {
	type args struct {
		task1 Task
		task2 Task
	}
	tests := []struct {
		name string
		args args
		want Task
	}{
		{
			name: "",
			args: args{
				task1: Task{
					TaskName2task: map[string]*Task{
						"task1": {
							TaskName2task: map[string]*Task{
								"task2": {
									TaskName2task: map[string]*Task{},
									Items: []*Item{
										{
											Content:  "item2.1",
											Comments: []string{},
											Level:    3,
										},
									},
									Level: 2,
									Name:  "task2",
								},
							},
							Items: []*Item{
								{
									Content:  "item1.1",
									Comments: []string{},
									Level:    1,
								},
							},
							Level:  0,
							Name:   "task1",
							Parent: nil,
						},
					},
					Items:  []*Item{},
					Level:  -1,
					Name:   "head",
					Parent: nil,
				},
				task2: Task{
					TaskName2task: map[string]*Task{
						"task1": {
							TaskName2task: map[string]*Task{
								"task2": {
									TaskName2task: map[string]*Task{},
									Items: []*Item{
										{
											Content:  "item2.2",
											Comments: []string{},
											Level:    3,
										},
									},
									Level: 2,
									Name:  "task2",
								},
							},
							Items: []*Item{
								{
									Content:  "item1.2",
									Comments: []string{},
									Level:    1,
								},
							},
							Level:  0,
							Name:   "task1",
							Parent: nil,
						},
					},
					Items:  []*Item{},
					Level:  -1,
					Name:   "head",
					Parent: nil,
				},
			},
			want: Task{
				TaskName2task: map[string]*Task{
					"task1": {
						TaskName2task: map[string]*Task{
							"task2": {
								Name:          "task2",
								TaskName2task: map[string]*Task{},
								Items: []*Item{
									{
										Content:  "item2.1",
										Comments: []string{},
										Level:    3,
									},
									{
										Content:  "item2.2",
										Comments: []string{},
										Level:    3,
									},
								},
							},
						},
						Items: []*Item{
							{
								Content:  "item1.1",
								Comments: []string{},
								Level:    1,
							},
							{
								Content:  "item1.2",
								Comments: []string{},
								Level:    1,
							},
						},
						Level:  0,
						Name:   "task1",
						Parent: nil,
					},
				},
				Items:  []*Item{},
				Level:  -1,
				Name:   "head",
				Parent: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MergeTasks(tt.args.task1, tt.args.task2); !IsTasksEqual(&got, &tt.want) {
				t.Errorf("MergeTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func IsTasksEqual(task1, task2 *Task) bool {
func IsTasksEqual(task1, task2 *Task) bool {
	if task1 == nil && task2 == nil {
		return true
	}
	if task1 == nil || task2 == nil {
		return false
	}
	if task1.Name != task2.Name {
		return false
	}
	if len(task1.Items) != len(task2.Items) {
		return false
	}
	for i := range task1.Items {
		if task1.Items[i].Content != task2.Items[i].Content {
			return false
		}
	}
	if len(task1.TaskName2task) != len(task2.TaskName2task) {
		return false
	}
	for k := range task1.TaskName2task {
		if !IsTasksEqual(task1.TaskName2task[k], task2.TaskName2task[k]) {
			return false
		}
	}
	return true
}

func TestParseTaskFromTodoFile(t *testing.T) {
	// TODO (xiangjun.ma) create unit test to write back and compare
	type args struct {
		todoFile string
	}
	tests := []struct {
		name    string
		args    args
		want    Task
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				// todoFile: "ut.todo",
				todoFile: "2023-01-13.todo",
			},
			want:    Task{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTaskFromTodoFile(tt.args.todoFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTaskFromTodoFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseTaskFromTodoFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseThenWrite(t *testing.T) {
	got, err := ParseTaskFromTodoFile("2023-01-13.todo")
	if err != nil {
		t.Errorf("ParseTaskFromTodoFile() error = %v", err)
		return
	}
	content := got.ToContent()
	// write to "ut_out.todo"
	err = ioutil.WriteFile("ut_out.todo", []byte(content), 0o644)
	if err != nil {
		t.Errorf("WriteFile() error = %v", err)
		return
	}
}
