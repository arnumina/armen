CREATE TABLE plugins (
	name VARCHAR(10) NOT NULL,
	config JSONB NOT NULL,
	CONSTRAINT plugins_pk PRIMARY KEY (name ASC),
	FAMILY "primary" (name, config)
)
