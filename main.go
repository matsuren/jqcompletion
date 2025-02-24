package main

import (
	"fmt"
	"os"
)

func main() {
	var jsonPath string
	if len(os.Args) == 1 {
		fmt.Println("Usage: jqcompletion <json_file_path>")
        os.Exit(0)
	} else {
		jsonPath = os.Args[1]
	}

	ui := NewJSONQueryUI()
	if err := ui.LoadJsonFile(jsonPath); err != nil {
		fmt.Printf("Error before application: %v\n", err)
	}

	if err := ui.Run(); err != nil {
		fmt.Printf("Error running application: %v\n", err)
	}
}
