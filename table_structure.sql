CREATE TABLE hierarchy (
  hierarchy_id VARCHAR NOT NULL,
  hierarchy_name VARCHAR,
  PRIMARY KEY( hierarchy_id )
);

CREATE TABLE hierarchy_entry (
  hierarchy_id VARCHAR NOT NULL,
  entry_code VARCHAR NOT NULL,
  parent_code VARCHAR,
  code_name VARCHAR,
  abbreviation VARCHAR,
  description VARCHAR,
  level INTEGER,
  PRIMARY KEY( hierarchy_id, entry_code ),
  FOREIGN KEY (hierarchy_id) REFERENCES hierarchy (hierarchy_id),
  FOREIGN KEY (hierarchy_id, parent_code) REFERENCES hierarchy_entry (hierarchy_id, entry_code)
);