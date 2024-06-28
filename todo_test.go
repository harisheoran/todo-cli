package todo_test

import (
	"cmd/todo"
	"os"
	"testing"
)

/*
func Assert(expected, actual int) {
	if expected != actual {
		t.Errorf("Expected %d but result is %d", expected, actual)
	}
}
*/

// testing Add fucntionality
func TestAdd(t *testing.T) {
	l := todo.List{}

	l.Add("Run 5 KM.")

	expectedRes := "Run 5 KM."

	result := l[0].Task

	if expectedRes != result {
		t.Errorf("Expected %s but result is %s", expectedRes, result)
	}
}

// testing complete fucntionality

func TestComplete(t *testing.T) {
	l := todo.List{}
	taskname := "Run 5 KM."
	l.Add(taskname)

	if l[0].Task != taskname {
		t.Errorf("Expected %s but result is %s", taskname, l[0].Task)
	}

	if l[0].Done {
		t.Errorf("Already task is %t", l[0].Done)
	}

	l.Complete(1)

	expectedRes := true

	result := l[0].Done

	if expectedRes != result {
		t.Errorf("Expected %t but result is %t", expectedRes, result)
	}

}

func TestDelete(t *testing.T) {
	l := todo.List{}
	l.Add("New task")
	l.Delete(1)

	expectedRes := 0
	result := len(l)

	if expectedRes != result {
		t.Errorf("Expected %d but result is %d", expectedRes, result)
	}

}

func TestSaveGet(t *testing.T) {
	list1 := todo.List{}
	list2 := todo.List{}
	task := "New Task"
	list1.Add(task)
	if task != list1[0].Task {
		t.Errorf("Expected %s but result is %s", task, list1[0].Task)
	}

	tempFile, err := os.CreateTemp("", "")

	if err != nil {
		t.Errorf("Error creatingtemp file")
	}
	defer os.Remove(tempFile.Name())

	if err := list1.Save(tempFile.Name()); err != nil {
		t.Fatalf("Error saving temp file %s", err)
	}

	if err := list2.Get(tempFile.Name()); err != nil {
		t.Fatalf("Error getting temp file %s", err)
	}

	if list1[0].Task != list2[0].Task {
		t.Errorf("Expected %s, actual result %s", list1[0].Task, list2[0].Task)
	}

}
