package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// item classs for task, not exposing it
type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type List []item

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
	if i <= 0 || i > len(*l) {
		return fmt.Errorf("No tasks in the list %d", i)
	}

	(*l)[i-1].Done = true
	(*l)[i-1].CompletedAt = time.Now()
	return nil

}

// Delete a task
func (l *List) Delete(i int) error {

	if i <= 0 || i > len(*l) {
		return fmt.Errorf("No tasks in the list %d", i)
	}

	*l = append((*l)[:i-1], (*l)[i:]...)
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
