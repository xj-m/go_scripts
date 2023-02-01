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

var (
	ItemStatusTodo = ItemStatus{
		Name:   "todo",
		Symbol: "❍",
	}
	ItemStatusDone = ItemStatus{
		Name:   "done",
		Symbol: "✔",
	}
	ItemStatusCanceled = ItemStatus{
		Name:   "canceled",
		Symbol: "✘",
	}
)

type ItemStatus struct {
	Name   string
	Symbol string
}

func ParseItemStatusFromLine(line string) (ItemStatus, error) {
	line = strings.TrimSpace(line)
	runes := []rune(line)
	if len(runes) == 0 {
		return ItemStatus{}, fmt.Errorf("line is empty")
	}
	switch runes[0] {
	case '❍':
		return ItemStatusTodo, nil
	case '✔':
		return ItemStatusDone, nil
	case '✘':
		return ItemStatusCanceled, nil
	default:
		return ItemStatus{}, fmt.Errorf("invalid item status symbol: %v", line[0])
	}
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
	Status        string
	Comments      []string
}

func (t *Task) RemoveTaskByNames(taskNames []string) {
	for _, taskName := range taskNames {
		delete(t.TaskName2task, taskName)
	}
	t.TaskNames = []string{}
	for taskName := range t.TaskName2task {
		t.TaskNames = append(t.TaskNames, taskName)
	}
	for _, task := range t.TaskName2task {
		task.RemoveTaskByNames(taskNames)
	}
}

func (t *Task) RemoveItem(itemToRemove *Item) {
	for i, item := range t.Items {
		if itemToRemove == item {
			t.Items = append(t.Items[:i], t.Items[i+1:]...)
			break
		}
	}
}

func (t *Task) SubTask(taskToRemove *Task) *Task {
	// copy t to ret
	ret := &Task{}
	mergo.Merge(ret, t)
	for taskName, subTask := range ret.TaskName2task {
		// check the taskName that also in taskToRemove
		if subTaskToRemove, ok := taskToRemove.TaskName2task[taskName]; ok {
			// scan items, remove the same one
			for _, item := range subTaskToRemove.Items {
				subTask.RemoveItem(item)
			}
			// check the subTask that also in taskToRemove
			subTask.SubTask(subTaskToRemove)
		}
	}
	return ret
}

func (t *Task) Filter(fs ...func(task Task) bool) {
	for _, task := range t.TaskName2task {
		task.Filter(fs...)
	}
	for taskName, task := range t.TaskName2task {
		for _, f := range fs {
			if !f(*task) {
				delete(t.TaskName2task, taskName)
				break
			}
		}
	}
	t.TaskNames = []string{}
	for taskName := range t.TaskName2task {
		t.TaskNames = append(t.TaskNames, taskName)
	}
}

func (t *Task) FilterItems(fs ...func(item Item) bool) {
	for _, task := range t.TaskName2task {
		task.FilterItems(fs...)
	}
	for i := 0; i < len(t.Items); i++ {
		for _, f := range fs {
			if !f(*t.Items[i]) {
				t.Items = append(t.Items[:i], t.Items[i+1:]...)
				i--
				break
			}
		}
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

	// if task is not root task, add task name and tags
	if t.Level != -1 {
		base = strings.Repeat("\t", t.Level)
		lines = []string{
			base + t.Name + ":" + t.TagK2v.ToContent(t.TagNames),
		}
	}
	// output comments
	for _, comment := range t.Comments {
		lines = append(lines, base+"\t"+comment)
	}

	// output items
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

	// join lines
	ret := strings.Join(lines, "\n")
	// replace tab with 4 spaces
	ret = strings.ReplaceAll(ret, "\t", "    ")
	return ret
}

func (t *Task) SortItemsByPriority() {
	for _, task := range t.TaskName2task {
		task.SortItemsByPriority()
	}
	sort.Sort(t.Items)
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
	Status     ItemStatus
	Content    string
	Comments   []string
	Level      int
	TagK2v     tagK2v
	TagNames   []string
	ParentTask *Task
	// TODO (xiangjun.ma) parser tags from content
	// TODO (xiangjun.ma) use status to indicate whether this item is done
}

func (item *Item) GetPriority() int {
	for k := range item.TagK2v {
		switch k {
		case "critical":
			return 100
		case "high":
			return 50
		case "low":
			return -50
		}
	}
	return 0
}

func (item *Item) ToContent() string {
	base := strings.Repeat("\t", item.Level)
	lines := []string{
		// base + item.Status + " " + item.Content + item.TagK2v.ToContent(item.TagNames),
		fmt.Sprintf("%s%s %s%s", base, item.Status.Symbol, item.Content, item.TagK2v.ToContent(item.TagNames)),
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
	status, err := ParseItemStatusFromLine(line)
	if err != nil {
		return err
	}
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

func (s Items) Len() int {
	return len(s)
}

func (s Items) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Items) Less(i, j int) bool {
	pI := s[i].GetPriority()
	pJ := s[j].GetPriority()
	return pI > pJ
}

var _ sort.Interface = Items{}

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
