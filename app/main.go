package main

import (
	"flag"
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
		defer f.Close()
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

	// Main
	jsonPath, err := getTempJsonPath()
	if err != nil {
		fmt.Print(err)
		os.Exit(0)
	}
	defer os.Remove(jsonPath)
	if len(args) == 1 {
		err := copyFile(args[0], jsonPath)
		if err != nil {
			fmt.Print(err)
			os.Exit(0)
		}
	} else if len(args) == 0 {
		cmd := editFileExecCmd(jsonPath)
		err := cmd.Run()
		if err != nil {
			fmt.Print(err)
			os.Exit(0)
		}
	} else {
		panic("Wrong arguments")
	}
	m := initializeModel()
	jsonData, err := LoadJsonFile(jsonPath)
	if err != nil {
		fmt.Print(err)
		os.Exit(0)
	}
	m.SetJsonData(jsonData)

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	if finalModel, ok := finalModel.(model); ok {
		fmt.Printf("jq '%v%v'", finalModel.queryHist, finalModel.queryEval)
	}
}
