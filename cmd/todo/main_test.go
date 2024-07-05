package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	binaryFile = "mytodo"
	filename   = ".todo.json"
)

func TestMain(m *testing.M) {
	fmt.Println("Building tool...")

	if runtime.GOOS == "windows" {
		binaryFile += ".exe"
	}

	cmd := exec.Command("touch", filename)
	lsErr := cmd.Run()
	if lsErr != nil {
		fmt.Println("THIS", lsErr)
		os.Exit(1)
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
	task := "Go to the GYM"
	dir, dirErr := os.Getwd()
	if dirErr != nil {
		t.Fatal(dirErr)
	}

	cmdPath := filepath.Join(dir, binaryFile)

	t.Run("Add new task", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-task", task)

		err := cmd.Run()

		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("List Task", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := task + "\n"

		if expected != string(output) {
			t.Errorf("Expected %q , result is %s", expected, string(output))
		}

	})

}
