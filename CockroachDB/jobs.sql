CREATE TABLE jobs (
	id UUID NOT NULL,
	name VARCHAR(30) NOT NULL,
	plugin VARCHAR(10) NOT NULL,
	type VARCHAR(30) NOT NULL,
	origin VARCHAR(30) NOT NULL,
	priority INT2 NOT NULL,
	key VARCHAR(30) NULL,
	emails VARCHAR(50) NULL,
	data JSONB NOT NULL,
	status VARCHAR(10) NOT NULL,
	attempt INT2 NOT NULL,
	wf_id UUID NULL,
	wf_failed BOOL NULL,
	reference TIMESTAMPTZ NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	run_after TIMESTAMPTZ NOT NULL,
	result VARCHAR(10) NULL,
	next_step VARCHAR(30) NULL,
	CONSTRAINT jobs_pk PRIMARY KEY (id ASC),
	CONSTRAINT jobs_fk FOREIGN KEY (wf_id) REFERENCES workflows(id) ON DELETE CASCADE,
	INDEX jobs_wf_id_idx (wf_id ASC),
	FAMILY "primary" (id, name, plugin, type, origin, priority, key, emails, data, status, attempt, wf_id, wf_failed, reference, created_at, run_after, result, next_step)
)
