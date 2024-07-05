package main

import (
	"cmd/todo"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var filename = "task.json"
	listFlag := flag.Bool("list", false, "List all the tasks")
	addFlag := flag.Bool("add", false, "Add a task")
	completeFlag := flag.Int("complete", -1, "Mark the task completed")
	deleteFlag := flag.Int("delete", -1, "Delete the task")
	verboseFlag := flag.Bool("details", false, "Detailed View of tasks")
	flag.Parse()

	if os.Getenv("TODO_FILENAME") != "" {
		filename = os.Getenv("TODO_FILENAME")
	}

	list := &todo.List{}

	err := list.Get(filename)
	if err != nil {
		log.Fatal(err)
	}

	switch {
	case *listFlag:

		if len(*list) == 0 {
			fmt.Println("No tasks")
		}

		// list all the tasks
		string := list.PrettyOutput()
		fmt.Println(string)
	case *addFlag:
		task, errGet := list.GetTask(os.Stdin, flag.Args()...)
		if err != nil {
			log.Fatal(errGet)
			os.Exit(1)
		}
		list.Add(task)
		err := list.Save(filename)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Task saved successfully.")

	case *completeFlag != -1 && *completeFlag >= 0:
		err := list.Complete(*completeFlag)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Task marked completed")
		errSave := list.Save(filename)
		if errSave != nil {
			log.Fatal(err)
		}
	case *deleteFlag > -1:
		err := list.Delete(*deleteFlag)
		if err != nil {
			log.Fatal(err)
		}
		errSave := list.Save(filename)
		if errSave != nil {
			log.Fatal(err)
		}
		fmt.Println("Task Deleted successfully.")
	case *verboseFlag:
		verboseData := list.VerboseOutput()
		fmt.Println(verboseData)

	default:
		fmt.Fprintln(os.Stderr, "Invalid Options")
		os.Exit(1)
	}
}
