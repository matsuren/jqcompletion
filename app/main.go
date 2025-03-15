package main

import (
	"fmt"
	"io"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
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

	jsonPath, err := getTempJsonPath()
	if err != nil {
		fmt.Printf("Got err: %v", err)
		os.Exit(0)
	}
	defer os.Remove(jsonPath)
	if len(os.Args) == 2 {
		err := copyFile(os.Args[1], jsonPath)
		if err != nil {
			fmt.Printf("Got err: %v", err)
			os.Exit(0)
		}
	}
	if len(os.Args) == 1 {
		cmd := editFileExecCmd(jsonPath)
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Got err: %v", err)
			os.Exit(0)
		}
	}
	m := initializeModelWithJsonFile(jsonPath)
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
