package task

import (
	"strings"

	"github.com/xj-m/go_scripts/file"
)

func TaskToFile(task Task, todoFile string) error {
	// FIXME (xiangjun.ma) write to file
	return nil
}

func ParseTaskFromTodoFile(todoFile string) (Task, error) {
	// read file content by line
	lines, err := file.ReadLines(todoFile)
	if err != nil {
		return Task{}, err
	}

	head := Task{
		TaskName2task: map[string]*Task{},
		Items:         []*Item{},
		Level:         -1,
		Name:          "head",
		Parent:        nil,
	}
	curTask := &head
	curItem := &Item{}
	for _, line := range lines {
		parseRes := parseLine(line)
		// curLevel, content, tagK2v, tagNames := getLevelContentTagsAndTagNames(line)
		runes := []rune(parseRes.Content)
		switch {
		case len(runes) == 0:
			// skip empty line
			continue
		case isTask(runes):
			// create new task and add to parent task
			parentTask := curTask.FindParentTask(parseRes.Level)
			newTask := Task{}
			err := newTask.ParseFromLine(line, parentTask)
			if err != nil {
				return Task{}, err
			}
			// if task already exists, merge two task and replace it
			if existTask, ok := parentTask.TaskName2task[newTask.Name]; ok {
				*existTask = MergeTasks(*existTask, newTask)
				curTask = existTask
			} else {
				parentTask.TaskName2task[newTask.Name] = &newTask
				parentTask.TaskNames = append(parentTask.TaskNames, newTask.Name)
				curTask = &newTask
			}
			curItem = nil
		case isItem(runes):
			// create item based on content
			newItem := Item{}
			err := newItem.ParseFromLine(line)
			if err != nil {
				return Task{}, err
			}
			if newItem.Content == "" {
				continue
			}
			// if newItem belongs to parent, then set curTask back to parent
			if newItem.Level < curTask.Level+1 {
				curTask = curTask.Parent
			}
			curTask.Items = append(curTask.Items, &newItem)
			curItem = &newItem
		default:
			// else treat as comment
			if len(curTask.TaskName2task) == 0 {
				curTask.Comments = append(curTask.Comments, parseRes.Content)
			} else {
				curItem.Comments = append(curItem.Comments, parseRes.Content)
			}
		}
	}
	return head, nil
}

func MergeTasks(srcTask, dstTask Task) Task {
	ret := dstTask
	// merge comments
	for _, comment := range srcTask.Comments {
		if !contains(ret.Comments, comment) {
			ret.Comments = append(ret.Comments, comment)
		}
	}

	// merge taskName2task
	for taskName, srcSubTask := range srcTask.TaskName2task {
		if dstSubTask, ok := dstTask.TaskName2task[taskName]; ok {
			mergedTask := MergeTasks(*srcSubTask, *dstSubTask)
			ret.TaskName2task[taskName] = &mergedTask
		} else {
			taskToAdd := *srcSubTask
			ret.TaskName2task[taskName] = &taskToAdd
			ret.TaskNames = append(ret.TaskNames, taskName)
		}
	}
	for _, itemI := range srcTask.Items {
		// if item is already in ret, skip
		if dstTask.Items.Contains(*itemI) {
			continue
		}
		ret.Items = append(ret.Items, itemI)
	}

	// merge TagK2V
	for k, v := range srcTask.TagK2v {
		if _, ok := ret.TagK2v[k]; !ok {
			ret.TagK2v[k] = v
		}
	}
	return ret
}

type ParseResult struct {
	Level    int
	Content  string
	TagK2v   map[string]string
	TagNames []string
}

func parseLine(line string) ParseResult {
	if len(line) > 3 && line[:2] == "  " && line[2] != ' ' {
		// special case for "Archive", where 2 spaces may means 1 tab
		line = "    " + line[2:]
	}
	// find and replace 4 spaces with 1 tab
	i := 0
	for ; i < len(line); i++ {
		if line[i] == ' ' && i+3 < len(line) && line[i:i+4] == "    " {
			line = line[:i] + "\t" + line[i+4:]
		} else {
			break
		}
	}
	// return number of leading tabs
	curLevel := 0
	for _, c := range line {
		if c == '\t' {
			curLevel++
		} else {
			break
		}
	}
	content := line[curLevel:]
	// trim leading and trailing spaces
	content = strings.TrimSpace(content)
	// if find @, treat as tag, following as tag key, inside the () as tag value
	tagKey2value := map[string]string{}
	tagNames := []string{}
	i = 0
	for i < len(content) {
		if content[i] == '@' {
			// find tag key
			j := i + 1
			hasValue := false
			tagValue := ""
			for ; j < len(content); j++ {
				if content[j] == '(' {
					hasValue = true
					break
				}
				if content[j] == ' ' {
					break
				}
			}
			tagKey := content[i+1 : j]
			tagNames = append(tagNames, tagKey)
			// find tag value
			if hasValue {
				k := j + 1
				for ; k < len(content); k++ {
					if content[k] == ')' {
						break
					}
				}
				tagValue = content[j+1 : k]
				tagKey2value[tagKey] = tagValue
				// remove tag from content
				content = content[:i] + content[k+1:]
			} else {
				tagKey2value[tagKey] = ""
				// remove tag from content
				content = content[:i] + content[j:]
			}
		}
		i++
	}
	content = strings.TrimSpace(content)
	return ParseResult{
		Level:    curLevel,
		Content:  content,
		TagK2v:   tagKey2value,
		TagNames: tagNames,
	}
}

func isTask(runes []rune) bool {
	return runes[len(runes)-1] == ':'
}

func isItem(runes []rune) bool {
	switch runes[0] {
	case '❍', '✔', '✘':
		return true
	default:
		return false
	}
}
