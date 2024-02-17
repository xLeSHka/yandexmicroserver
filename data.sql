CREATE TABLE expressions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    expression VARCHAR(256),
    expression_status VARCHAR(256),
    created_at TIMESTAMP
);
CREATE TABLE operations (
    operation VARCHAR(1),
    execution_time_by_milliseconds INTEGER DEFAULT 2000
);
CREATE TABLE agents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    status_code VARCHAR(256)
);
INSERT INTO 
		operations (execution_time_by_milliseconds) 
		VALUES (1000) WHERE operation = '+';
INSERT INTO 
		public.operations (operation,execution_time_by_milliseconds) 
		VALUES (1000);
INSERT INTO 
		public.operations (operation,execution_time_by_milliseconds) 
		VALUES (1000);
INSERT INTO 
		public.operations (operation,execution_time_by_milliseconds) 
		VALUES (1000);