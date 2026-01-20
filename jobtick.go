package jobtick

import (
	"database/sql"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// JobTick is the main library struct.
type JobTick struct {
	db *sql.DB
}

// New creates a new JobTick instance with a database connection.
func New(db *sql.DB) *JobTick {
	return &JobTick{db: db}
}

// CreateProcess inserts a new process into the database.
// The delay argument can be a time.Time (absolute start time) or an int64 (delay in seconds from now).
func (jt *JobTick) CreateProcess(proc Process, delay interface{}) error {
	var startTime time.Time

	switch d := delay.(type) {
	case time.Time:
		startTime = d
	case int64:
		startTime = time.Now().Add(time.Duration(d) * time.Second)
	case int:
		startTime = time.Now().Add(time.Duration(d) * time.Second)
	default:
		startTime = time.Now().Add(1 * time.Minute) // default 1 minute delay
	}

	return jt.Create(proc.Name, proc.Type, proc.Queue, proc.File, proc.Function, proc.Data, startTime)
}

// Execute fetches the next pending process for a queue and runs it.
func (jt *JobTick) Execute(queue string) error {
	// 1. Fetch and mark as assigned
	proc, err := jt.GetNextForExecution(queue)
	if err != nil {
		return err
	}

	if proc == nil {
		return nil
	}

	return jt.ExecuteProcess(proc)
}

// ExecuteProcess runs the process by executing its associated file.
func (jt *JobTick) ExecuteProcess(proc *Process) error {
	// 1. Ensure it's marked as assigned if it wasn't already (e.g. for forced execution)
	if proc.Status != StatusAssigned {
		err := jt.UpdateStatus(proc.ID, StatusAssigned)
		if err != nil {
			return err
		}
		proc.Status = StatusAssigned
	}

	// 2. Execute the file with processStart, function, and data
	cmd := exec.Command(proc.File, "processStart", proc.Function, string(proc.Data))
	output, err := cmd.Output()

	// 3. Update status based on result
	status := StatusExecuted
	if err != nil {
		status = StatusFailed
	}

	err = jt.UpdateStatus(proc.ID, status)
	if err != nil {
		return err
	}

	// 4. Handle recurring process
	if proc.Type == Recurring && status == StatusExecuted {
		delayStr := strings.TrimSpace(string(output))
		if delay, parseErr := strconv.ParseInt(delayStr, 10, 64); parseErr == nil {
			return jt.CreateProcess(*proc, delay)
		}
	}

	return nil
}
