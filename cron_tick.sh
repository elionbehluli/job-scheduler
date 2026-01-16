#!/bin/bash

# Get the absolute path of the project directory
PROJECT_DIR="/Users/elion/Desktop/job/job-scheduler"

# Change to the project directory
cd "$PROJECT_DIR" || exit

# Run the job-scheduler with the queue command
# Format: jobtick queue <queue_name>
go run cmd/jobtick/main.go "queue" "default"

# To add this to your crontab, run 'crontab -e' and add the following line:
# * * * * * /Users/elion/Desktop/job/job-scheduler/cron_tick.sh >> /Users/elion/Desktop/job/job-scheduler/cron.log 2>&1
