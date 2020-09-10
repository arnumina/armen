CREATE TABLE events (
	name VARCHAR(50) NOT NULL,
	disabled BOOL NOT NULL,
	after VARCHAR(10) NULL,
	repeat VARCHAR(50) NULL,
	data JSONB NULL,
	CONSTRAINT scheduler_pk PRIMARY KEY (name ASC),
	FAMILY "primary" (name, disabled, after, repeat, data)
)
