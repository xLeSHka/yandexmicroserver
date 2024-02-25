CREATE TABLE expressions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    expression VARCHAR(256) NOT NULL,
    expression_status VARCHAR(256) NOT NULL,
    created_at TIMESTAMP NOT NULL,
	completed_at TIMESTAMP NOT NULL
);
CREATE TABLE operations (
    operation VARCHAR(1) NOT NULL,
    execution_time_by_milliseconds INTEGER NOT NULL DEFAULT 2000
);
CREATE TABLE agents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	agent_address VARCHAR(100) NOT NULL,
    status_code VARCHAR(256) NOT NULL,
	last_heartbeat TIMESTAMP NOT NULL
);
UPDATE
	public.operations SET execution_time_by_milliseconds
	= $1 WHERE operation = $2
UPDATE public.expressions SET (expression,expression_status,completed_at) 
	= ($1,$2,$3) WHERE id = $4;		