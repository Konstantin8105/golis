package golis_test

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TODO add instal lis in travis

func TestTodo(t *testing.T) {
	// Show all to do`s in comment code
	source, err := filepath.Glob(fmt.Sprintf("./%s", "*.go"))
	if err != nil {
		t.Fatal(err)
	}

	var amount int

	for i := range source {
		t.Run(source[i], func(t *testing.T) {
			file, err := os.Open(source[i])
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = file.Close() }()

			pos := 0
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				pos++
				index := strings.Index(line, "//")
				if index < 0 {
					continue
				}
				if !strings.Contains(strings.ToUpper(line), "TODO") {
					continue
				}
				t.Logf("%13s:%-4d %s", source[i], pos, line[index:])
				amount++
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
		})
	}
	if amount > 0 {
		t.Logf("Amount TODO: %d", amount)
	}
}

func TestFmt(t *testing.T) {
	// Show all fmt`s in comments code
	source, err := filepath.Glob(fmt.Sprintf("./%s", "*.go"))
	if err != nil {
		t.Fatal(err)
	}

	var amount int

	for i := range source {
		t.Run(source[i], func(t *testing.T) {
			file, err := os.Open(source[i])
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = file.Close() }()

			pos := 1
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				pos++
				index := strings.Index(line, "//")
				if index < 0 {
					continue
				}
				if !strings.Contains(line, "fmt.") {
					continue
				}
				t.Logf("%d %s", pos, line[index:])
				amount++
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
		})
	}
	if amount > 0 {
		t.Logf("Amount commented fmt: %d", amount)
	}
}

func TestDebug(t *testing.T) {
	// Show all debug information in code
	source, err := filepath.Glob(fmt.Sprintf("./%s", "*.go"))
	if err != nil {
		t.Fatal(err)
	}

	for i := range source {
		t.Run(source[i], func(t *testing.T) {
			file, err := os.Open(source[i])
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = file.Close() }()

			pos := 1
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				pos++
				if !strings.Contains(line, "fmt"+"."+"Print") {
					continue
				}
				t.Errorf("Debug line: %d in file %s", pos, source[i])
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
		})
	}
}
