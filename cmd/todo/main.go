package main

import (
	"cmd/todo"
	"fmt"
	"os"
	"strings"
)

func main() {
	const filename = "task.json"

	list := &todo.List{}

	err := list.Get(filename)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch len(os.Args) {
	case 1:
		for i, value := range *list {
			fmt.Println("%d . %s", i, value.Task)
		}
	default:
		task := strings.Join(os.Args[1:], " ")
		list.Add(task)
		err := list.Save(filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	}

}
