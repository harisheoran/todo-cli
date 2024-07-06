package todo

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type List []item

type Stringer interface {
	PrettyOutput() string
	VerboseOutput() string
}

// item classs for task, not exposing it
type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

func (list *List) VerboseOutput() string {
	verbose := ""
	for i, value := range *list {
		if value.CompletedAt.Format("01-02-2006 15:04:05") == "01-01-0001 00:00:00" {
			verbose += fmt.Sprintf("%d %s CREATED:%s COMPLETED:%s \n", i, value.Task, value.CreatedAt.Format("01-02-2006 15:04:05"), "In Progress")
		} else {
			verbose += fmt.Sprintf("%d %s CREATED:%s COMPLETED:%s \n", i, value.Task, value.CreatedAt.Format("01-02-2006 15:04:05"), value.CompletedAt.Format("01-02-2006 15:04:05"))
		}
	}
	return verbose
}

func (list *List) PrettyOutput() string {
	formattedOutput := ""

	for i, value := range *list {
		prefix := "  "
		if value.Done {
			prefix = "X "
		}
		formattedOutput += fmt.Sprintf("%s%d: %s\n", prefix, i, value.Task)
	}

	return formattedOutput
}

func (l *List) GetTask(r io.Reader, args ...string) ([]string, error) {
	task := []string{}

	if len(args) > 0 {
		task = append(task, strings.Join(args, " "))
		return task, nil
	}

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		task = append(task, s.Text())
		if len(line) == 0 {
			break
		}
	}

	if err := s.Err(); err != nil {
		return task, err
	}

	if len(s.Text()) == 0 {
		return task, fmt.Errorf("Task Cannot be blank")
	}

	return task, nil
}

// Add the task
func (l *List) Add(task string) {
	mytask := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}
	*l = append(*l, mytask)
}

// Mark the task as done
func (l *List) Complete(i int) error {
	// Check if List have tasks or not
	if i < 0 || i > len(*l) {
		return fmt.Errorf("No tasks in the list %d", i)
	}

	(*l)[i].Done = true
	(*l)[i].CompletedAt = time.Now()
	return nil
}

// Delete a task
func (l *List) Delete(i int) error {

	if i < 0 || i > len(*l) {
		return fmt.Errorf("No tasks in the list %d", i)
	}

	*l = append((*l)[:i], (*l)[i+1:]...)
	return nil
}

func (l *List) Save(filename string) error {
	_, err := os.Stat(filename)

	if errors.Is(err, os.ErrExist) {
		os.Create(filename)
	}

	data, errMarshal := json.Marshal(*l)

	if errMarshal != nil {
		return fmt.Errorf("Error Mashaling the task into JSON %d", errMarshal)
	}

	errWrite := os.WriteFile(filename, data, 0644)
	if errWrite != nil {
		return fmt.Errorf("Error saving the task %d", errWrite)
	}

	return nil
}

// Get a task
func (l *List) Get(filename string) error {
	data, err := os.ReadFile(filename)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if len(data) == 0 {
		return fmt.Errorf("Not any task saved")
	}

	return json.Unmarshal(data, l)
}
