package task

import (
	"fmt"
	"sort"
	"strings"

	"github.com/imdario/mergo"
)

type tagK2v map[string]string

func (t tagK2v) ToContent(names []string) string {
	ret := ""
	for _, name := range names {
		v := t[name]
		if v == "" {
			ret += fmt.Sprintf(" @%s", name)
			continue
		}
		ret += fmt.Sprintf(" @%s(%s)", name, v)
	}
	return ret
}

type Task struct {
	TaskName2task map[string]*Task
	TaskNames     []string
	Items         Items
	Level         int
	Name          string
	Parent        *Task
	TagK2v        tagK2v
	TagNames      []string
}

func (t *Task) RemoveTaskByNames(taskNames []string) {
	for _, taskName := range taskNames {
		delete(t.TaskName2task, taskName)
	}
	t.TaskNames = []string{}
	for taskName := range t.TaskName2task {
		t.TaskNames = append(t.TaskNames, taskName)
	}
}

func (t *Task) ParseFromLine(line string, parentTask *Task) error {
	parseRes := parseLine(line)
	taskName := parseRes.Content[:len(parseRes.Content)-1]
	parsedTask := &Task{
		TaskName2task: map[string]*Task{},
		Items:         []*Item{},
		Level:         parseRes.Level,
		Name:          taskName,
		Parent:        parentTask,
		TagK2v:        parseRes.TagK2v,
		TagNames:      parseRes.TagNames,
	}
	mergo.Merge(t, parsedTask)
	return nil
}

func (t *Task) FindParentTask(level int) *Task {
	if t.Level == level-1 {
		return t
	}
	return t.Parent.FindParentTask(level)
}

func (t *Task) ToContent() string {
	base, lines := "", []string{}
	if t.Level != -1 {
		base = strings.Repeat("\t", t.Level)
		lines = []string{
			base + t.Name + ":" + t.TagK2v.ToContent(t.TagNames),
		}
	}
	for _, item := range t.Items {
		lines = append(lines, item.ToContent())
	}
	ts := make(Tasks, 0, len(t.TaskName2task))
	for _, taskName := range t.TaskNames {
		ts = append(ts, t.TaskName2task[taskName])
	}
	for i, task := range ts {
		lines = append(lines, task.ToContent())
		// if is last sub-task, add a blank line
		if i != len(t.TaskName2task)-1 {
			lines = append(lines, "")
		}
	}
	if len(t.Items) == 0 && len(t.TaskName2task) == 0 {
		lines = append(lines, base+"\t❍ ")
	}
	ret := strings.Join(lines, "\n")
	// replace tab with 4 spaces
	ret = strings.ReplaceAll(ret, "\t", "    ")
	return ret
}

type Tasks []*Task

func (s Tasks) Len() int {
	return len(s)
}

func (s Tasks) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Tasks) Less(i, j int) bool {
	pI, okI := TaskName2Priority[s[i].Name]
	if !okI {
		pI = 0
	}
	pJ, okJ := TaskName2Priority[s[j].Name]
	if !okJ {
		pJ = 0
	}
	if pI == pJ {
		return s[i].Name < s[j].Name
	}
	return pI > pJ
}

var _ sort.Interface = Tasks{}

var TaskName2Priority = map[string]int{
	"routine": 100,
	"misc":    -50,
	"inbox":   -100,
}

type Item struct {
	Status   string
	Content  string
	Comments []string
	Level    int
	TagK2v   tagK2v
	TagNames []string
	// TODO (xiangjun.ma) parser tags from content
	// TODO (xiangjun.ma) use status to indicate whether this item is done
}

func (item *Item) ToContent() string {
	base := strings.Repeat("\t", item.Level)
	lines := []string{
		// base + item.Status + " " + item.Content + item.TagK2v.ToContent(item.TagNames),
		fmt.Sprintf("%s%s %s%s", base, item.Status, item.Content, item.TagK2v.ToContent(item.TagNames)),
	}
	for _, comment := range item.Comments {
		lines = append(lines, base+"\t"+comment)
	}

	ret := strings.Join(lines, "\n")
	ret = strings.ReplaceAll(ret, "\t", "    ")
	return ret
}

func (item *Item) ParseFromLine(line string) error {
	parseRes := parseLine(line)
	status := string([]rune(parseRes.Content)[:1])
	content := strings.Trim(string([]rune(parseRes.Content)[1:]), " ")
	parsedItem := &Item{
		Status:   status,
		Content:  content,
		Comments: []string{},
		Level:    parseRes.Level,
		TagK2v:   parseRes.TagK2v,
		TagNames: parseRes.TagNames,
	}
	mergo.Merge(item, parsedItem)
	return nil
}

func GetValues[T any](m map[string]T) []T {
	var values []T
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

type Items []*Item

func (s Items) Contains(item Item) bool {
	for _, i := range s {
		if i.Content == item.Content && i.Level == item.Level && isEq(i.Comments, item.Comments) {
			return true
		}
	}
	return false
}

func isEq[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
