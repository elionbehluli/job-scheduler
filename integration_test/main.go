package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"jobtick"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
)

func main() {
	// 1. Database Connection
	connStr := "postgres://localhost:5432/jobtick_db?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		return
	}
	defer db.Close()

	jt := jobtick.New(db)

	// 2. Get the absolute path of the processor executable
	cwd, _ := os.Getwd()
	processorPath := filepath.Join(cwd, "processor")

	// 3. Create a recurring process
	proc := jobtick.Process{
		Name:     "Integration Test Job",
		Type:     jobtick.Recurring,
		Queue:    "test_queue",
		File:     processorPath,
		Function: "recurring_task",
		Data:     json.RawMessage(`{"msg": "Hello JobTick"}`),
	}

	fmt.Println("Creating recurring process...")
	err = jt.CreateProcess(proc, 0) // immediate start (actually scheduled for 0s from now)
	if err != nil {
		fmt.Printf("Error creating process: %v\n", err)
		return
	}

	// 4. Execute the queue
	fmt.Println("Executing 'test_queue'...")
	err = jt.Execute("test_queue")
	if err != nil {
		fmt.Printf("Error executing queue: %v\n", err)
		return
	}

	fmt.Println("Integration test complete. Check the database for a new scheduled run of this job.")
}
