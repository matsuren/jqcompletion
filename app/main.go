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
		fmt.Print(err)
		os.Exit(0)
	}
	defer os.Remove(jsonPath)
	if len(os.Args) == 2 {
		err := copyFile(os.Args[1], jsonPath)
		if err != nil {
			fmt.Print(err)
			os.Exit(0)
		}
	}
	if len(os.Args) == 1 {
		cmd := editFileExecCmd(jsonPath)
		err := cmd.Run()
		if err != nil {
			fmt.Print(err)
			os.Exit(0)
		}
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
