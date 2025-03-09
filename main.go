package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if os.Getenv("DEBUGLOG") != "" {
		f, err := tea.LogToFile("debug.log", "main")
		if err != nil {
			fmt.Println("Couldn't open a file for logging:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	var jsonPath string
	if len(os.Args) == 1 {
		fmt.Println("Usage: jqcompletion <json_file_path>")
	       os.Exit(0)
	} else {
		jsonPath = os.Args[1]
	}

	m := initializeModelWithJsonFile(jsonPath)
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
