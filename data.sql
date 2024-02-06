CREATE TABLE expressions (
    id INTEGER,
    expression VARCHAR(100),
    expression_status string,
    created_at TIMESTAMP,
    completed_at TIMESTAMP,
);
CREATE TABLE operations (
    operation VARCHAR(1),
    ExecuteTimeInSeconds integer,
)
CREATE TABLE deamons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    deamon_status INTEGER,
)