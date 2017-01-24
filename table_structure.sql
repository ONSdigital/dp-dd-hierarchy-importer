CREATE TABLE dimension (
  dimension_id VARCHAR NOT NULL,
  dimension_name VARCHAR,
  dimension_type VARCHAR,
  PRIMARY KEY( dimension_id )
);

CREATE TABLE dimension_level_type (
  type_id VARCHAR NOT NULL,
  type_name VARCHAR,
  type_level INTEGER,
  PRIMARY KEY( type_id )
);

CREATE TABLE dimension_value (
  dimension_id VARCHAR NOT NULL,
  value_code VARCHAR NOT NULL,
  parent_code VARCHAR,
  value_name VARCHAR,
  level_type VARCHAR,
  display_order INTEGER,
  PRIMARY KEY( dimension_id, value_code ),
  FOREIGN KEY (dimension_id) REFERENCES dimension (dimension_id),
  FOREIGN KEY (level_type) REFERENCES dimension_level_type (type_id),
  FOREIGN KEY (dimension_id, parent_code) REFERENCES dimension_value (dimension_id, value_code)
);