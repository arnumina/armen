CREATE TABLE armen (
	id UUID NOT NULL,
	host VARCHAR(30) NOT NULL,
	port INT4 NOT NULL,
	started_at TIMESTAMPTZ NOT NULL,
	CONSTRAINT armen_pk PRIMARY KEY (id ASC),
	UNIQUE INDEX armen_un (host ASC, port ASC),
	FAMILY "primary" (id, host, port, started_at)
)
