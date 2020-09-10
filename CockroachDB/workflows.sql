CREATE TABLE workflows (
	id UUID NOT NULL,
	name VARCHAR(30) NOT NULL,
	title VARCHAR(100) NOT NULL,
	origin VARCHAR(30) NOT NULL,
	priority INT2 NOT NULL,
	first_step VARCHAR(30) NOT NULL,
	steps JSONB NOT NULL,
	emails VARCHAR(50) NULL,
	data JSONB NOT NULL,
	status VARCHAR(10) NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	finished_at TIMESTAMPTZ NULL,
	CONSTRAINT workflows_pk PRIMARY KEY (id ASC),
	FAMILY "primary" (id, name, title, origin, priority, first_step, steps, emails, data, status, created_at, finished_at)
)
