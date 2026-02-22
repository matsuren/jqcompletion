package main

import (
	"flag"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if os.Getenv("DEBUGLOG") != "" {
		f, err := tea.LogToFile("debug.log", "queryview")
		if err != nil {
			fmt.Println("Couldn't open a file for logging:", err)
			os.Exit(1)
		}
		defer func() { _ = f.Close() }()
	} else {
		log.SetOutput(io.Discard)
		log.SetFlags(0)
	}

	// Parse arguments
	versionFlag := flag.Bool("v", false, "Print version")
	helpFlag := flag.Bool("h", false, "Print help")
	flag.Parse()
	args := flag.Args()

	//
	if *versionFlag {
		fmt.Printf("Version: %s, Commit: %s, Build Date: %s\n", version, commit, date)
		os.Exit(0)
	}
	//
	if *helpFlag || len(args) >= 2 {
		fmt.Println("Usage: jqcompletion [<filename>]")
		fmt.Printf("Version [-v]: %s, Commit: %s, Build Date: %s\n", version, commit, date)
		os.Exit(0)
	}

	// Detect piped stdin
	stat, _ := os.Stdin.Stat()
	stdinIsPipe := (stat.Mode() & os.ModeCharDevice) == 0

	// Main
	jsonPath, err := getTempJsonPath()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	defer func() { _ = os.Remove(jsonPath) }()
	if stdinIsPipe && len(args) == 1 {
		fmt.Fprintln(os.Stderr, "Error: ambiguous input — both stdin pipe and filename provided")
		os.Exit(1)
	} else if stdinIsPipe && len(args) == 0 {
		if err := readStdinToFile(os.Stdin, jsonPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			os.Exit(1)
		}
	} else if len(args) == 1 {
		err := copyFile(args[0], jsonPath)
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
	} else if len(args) == 0 {
		cmd := editFileExecCmd(jsonPath)
		err := cmd.Run()
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
	} else {
		panic("Wrong arguments")
	}
	m := initializeModel()
	jsonData, isSlurp, err := LoadJsonFile(jsonPath)
	if errors.Is(err, errEmptyInput) {
		return
	}
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	m.isSlurp = isSlurp
	m.SetJsonData(jsonData)

	// When stdin was piped, reopen /dev/tty for terminal input
	opts := []tea.ProgramOption{tea.WithAltScreen(), tea.WithMouseCellMotion()}
	if stdinIsPipe {
		ttyFile, err := os.Open("/dev/tty")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: could not open /dev/tty for terminal input")
			os.Exit(1)
		}
		defer func() { _ = ttyFile.Close() }()
		opts = append(opts, tea.WithInput(ttyFile))
	}
	p := tea.NewProgram(m, opts...)
	finalModel, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	if finalModel, ok := finalModel.(model); ok {
		if finalModel.isSlurp {
			fmt.Printf("jq -s '%v%v'\n", finalModel.queryHist, finalModel.queryEval)
		} else {
			fmt.Printf("jq '%v%v'\n", finalModel.queryHist, finalModel.queryEval)
		}
	}
}

func readStdinToFile(r io.Reader, dstPath string) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("reading stdin: %w", err)
	}
	if err := os.WriteFile(dstPath, data, 0o600); err != nil {
		return fmt.Errorf("writing temp file: %w", err)
	}
	return nil
}
