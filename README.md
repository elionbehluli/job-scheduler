# JobTick

A minimalist background job library for Go and Postgres.

## Integration

To use JobTick in your project:

1. **Install**:
   ```bash
   go get github.com/yourusername/jobtick
   ```

2. **Database Setup**:
   Run the SQL migration found in `migrations/0001_create_jobs_table.up.sql` on your Postgres database.

3. **Initialize**:
   ```go
   import (
       "jobtick"
       _ "github.com/lib/pq"
   )

   db, _ := sql.Open("postgres", "your_connection_string")
   jt := jobtick.New(db)
   ```

## Usage

### 1. Create a Process
```go
proc := jobtick.Process{
    Name:     "process_invoice",
    Type:     jobtick.Call, // or jobtick.Recurring
    Queue:    "billing",
    File:     "/path/to/processor",
    Function: "GeneratePDF",
    Data:     json.RawMessage(`{"id": 123}`),
}

// Start with a 10 second delay
err := jt.CreateProcess(proc, 10) 
```

### 2. The Runner (Cron)
Your worker file (e.g. `worker.go`) should look like this:

```go
func main() {
    // ... init jt ...
    // Executes the next pending job in the "billing" queue
    err := jt.Execute("billing")
}
```

### 3. Dynamic Execution
Each `File` specified in a process must handle the `processStart` argument:

```go
// /path/to/processor
func main() {
    if os.Args[1] == "processStart" {
        function := os.Args[2]
        data := os.Args[3]
        // ... run your logic ...
        
        // If Recurring, print seconds to wait for next run
        fmt.Println("60") 
    }
}
```
