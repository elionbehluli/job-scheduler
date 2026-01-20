package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Not enough arguments")
		os.Exit(1)
	}

	command := os.Args[1]
	function := os.Args[2]
	data := os.Args[3]

	if command == "processStart" {
		fmt.Fprintf(os.Stderr, "Worker starting function: %s with data: %s\n", function, data)

		// Simulate logic
		if function == "recurring_task" {
			// For recurring tasks, we print the delay for the next run (e.g. 5 seconds)
			fmt.Println("5")
		} else {
			fmt.Println("success")
		}
	} else {
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}
