package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var (
	binaryFile = "mytodo"
	filename   = "test.json"
)

func TestMain(m *testing.M) {
	fmt.Println("Building tool...")

	if runtime.GOOS == "windows" {
		binaryFile += ".exe"
	}

	command := exec.Command("go", "build", "-o", binaryFile)
	err := command.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot build tool %s:%s", binaryFile, err)
		os.Exit(1)
	}

	fmt.Println("Running tests...")

	result := m.Run()

	fmt.Println("Cleaning up...")

	os.Remove(binaryFile)
	os.Remove(filename)

	os.Exit(result)
}

func TestCLI(t *testing.T) {
	task := "Code everyday."
	task2 := "Make projects."

	dir, dirErr := os.Getwd()
	if dirErr != nil {
		t.Fatal(dirErr)
	}

	cmdPath := filepath.Join(dir, binaryFile)

	t.Run("Add new task", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add", task)

		err := cmd.Run()

		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Add a new task from STDIN", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add")
		cmdStdIn, err := cmd.StdinPipe()

		if err != nil {
			t.Fatal(err)
		}

		io.WriteString(cmdStdIn, task2)
		cmdStdIn.Close()

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("List Task", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprintf("  0: %s\n  1: %s\n", task, task2)

		result := strings.TrimSuffix(string(output), "\n")
		if expected != result {
			t.Errorf("Expected %q , result is %s", expected, result)
		}

	})

	t.Run("Complete Task", func(t *testing.T) {
		cmdComplete := exec.Command(cmdPath, "-complete=0")
		errComplete := cmdComplete.Run()
		if errComplete != nil {
			t.Fatal(errComplete)
		}

		cmd := exec.Command(cmdPath, "-list")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprintf("X 0: %s\n  1: %s\n", task, task2)

		result := strings.TrimSuffix(string(output), "\n")
		if expected != result {
			t.Errorf("Expected %q , result is %s", expected, result)
		}

	})

}
