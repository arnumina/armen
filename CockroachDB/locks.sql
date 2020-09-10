CREATE TABLE locks (
	name VARCHAR(20) NOT NULL,
	expiration_datetime TIMESTAMPTZ NULL DEFAULT NULL,
	owner UUID NULL DEFAULT NULL,
	CONSTRAINT locks_pk PRIMARY KEY (name ASC),
	FAMILY "primary" (name, expiration_datetime, owner)
)
