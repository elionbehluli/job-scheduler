package main

import (
	"database/sql"
	"fmt"
	"jobtick"
	"os"
	"time"
)

func main() {
	// PHP-style argc and argv
	argc := len(os.Args)
	argv := os.Args

	if argc < 3 {
		fmt.Println("Usage: jobtick <queue|force> <id|name>")
		return
	}

	command := argv[1]
	target := argv[2]

	// Initialize your database and JobTick
	// This is a placeholder - user must provide real DB connection
	var db *sql.DB
	jt := jobtick.New(db)

	var proc *jobtick.Process
	var err error

	switch command {
	case "queue":
		proc, err = jt.GetNextForExecution(target)
	case "force":
		proc, err = jt.FromID(target)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		return
	}

	if err != nil {
		fmt.Printf("Error fetching process: %v\n", err)
		return
	}

	if proc == nil {
		fmt.Printf("No process found for %s: %s\n", command, target)
		return
	}

	fmt.Printf("Executing process: %s (ID: %s)\n", proc.Name, proc.ID)

	err = jt.ExecuteProcess(proc)
	if err != nil {
		fmt.Printf("Execution failed: %v\n", err)
		return
	}

	fmt.Printf("[%s] Job complete\n", time.Now().Format(time.RFC822))
}
