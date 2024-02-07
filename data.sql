CREATE TABLE expressions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    expression VARCHAR(100),
    expression_status string,
    created_at TIMESTAMP,
    completed_at TIMESTAMP,
    execution_time_by_milliseconds INTEGER,
);
CREATE TABLE operations (
    operation VARCHAR(1),
    ExecuteTimeByMilliseconds INTEGER DEFAULT 2000,
)
CREATE TABLE deamons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    deamon_status INTEGER,
)

		
	

