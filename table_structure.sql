CREATE TABLE hierarchy (
  hierarchy_id VARCHAR NOT NULL,
  hierarchy_name VARCHAR,
  hierarchy_type VARCHAR,
  PRIMARY KEY( hierarchy_id )
);

CREATE TABLE hierarchy_level_type (
  type_id VARCHAR NOT NULL,
  type_name VARCHAR,
  type_level INTEGER,
  PRIMARY KEY( type_id )
);

CREATE TABLE hierarchy_entry (
  hierarchy_id VARCHAR NOT NULL,
  value_code VARCHAR NOT NULL,
  parent_code VARCHAR,
  value_name VARCHAR,
  level_type VARCHAR,
  display_order INTEGER,
  PRIMARY KEY( hierarchy_id, value_code ),
  FOREIGN KEY (hierarchy_id) REFERENCES hierarchy (hierarchy_id),
  FOREIGN KEY (level_type) REFERENCES hierarchy_level_type (type_id),
  FOREIGN KEY (hierarchy_id, parent_code) REFERENCES hierarchy_entry (hierarchy_id, value_code)
);