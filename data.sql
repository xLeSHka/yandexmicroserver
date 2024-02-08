CREATE TABLE database.expressions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    expression VARCHAR(100),
    expression_status VARCHAR(100),
    created_at TIMESTAMP
);
CREATE TABLE database.operations (
    operation VARCHAR(1),
    execution_time_by_milliseconds INTEGER DEFAULT 2000
);

