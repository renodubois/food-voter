CREATE TABLE poll (
	id INTEGER PRIMARY KEY ASC,
	options TEXT NOT NULL,
	creator_slug TEXT NOT NULL,
	voter_slug TEXT NOT NULL
);

CREATE TABLE result (
	id INTEGER PRIMARY KEY ASC,
	poll_id INTEGER NOT NULL,
	results TEXT NOT NULL,
	FOREIGN KEY (poll_id) REFERENCES poll(id)
);
