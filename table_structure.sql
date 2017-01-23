CREATE TABLE hierarchy (
  hierarchy_id VARCHAR NOT NULL,
  hierarchy_name VARCHAR,
  hierarchy_type VARCHAR,
  PRIMARY KEY( hierarchy_id )
);

CREATE TABLE hierarchy_level_type (
  id VARCHAR NOT NULL,
  name VARCHAR,
  level INTEGER,
  PRIMARY KEY( id )
);

CREATE TABLE hierarchy_entry (
  hierarchy_id VARCHAR NOT NULL,
  entry_code VARCHAR NOT NULL,
  parent_code VARCHAR,
  name VARCHAR,
  level_type VARCHAR,
  display_order INTEGER,
  PRIMARY KEY( hierarchy_id, entry_code ),
  FOREIGN KEY (hierarchy_id) REFERENCES hierarchy (hierarchy_id),
  FOREIGN KEY (level_type) REFERENCES hierarchy_level_type (id),
  FOREIGN KEY (hierarchy_id, parent_code) REFERENCES hierarchy_entry (hierarchy_id, entry_code)
);