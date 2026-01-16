CREATE TABLE IF NOT EXISTS processes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    type TEXT NOT NULL, -- 'single call', 'recurring'
    queue TEXT NOT NULL,
    file TEXT,
    function TEXT,
    data JSONB,
    start_time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    status TEXT NOT NULL DEFAULT 'created', -- 'created', 'assigned', 'executed'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_processes_status_start_time ON processes (status, start_time);
CREATE INDEX IF NOT EXISTS idx_processes_queue ON processes (queue);
