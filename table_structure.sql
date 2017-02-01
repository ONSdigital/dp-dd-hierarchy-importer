
CREATE TABLE "hierarchy"
(
  id varchar PRIMARY KEY NOT NULL,
  name varchar,
  type varchar
);

CREATE TABLE "hierarchy_level_type"
(
  id varchar PRIMARY KEY NOT NULL,
  level int,
  name varchar
);

CREATE TABLE "hierarchy_entry"
(
  id varchar PRIMARY KEY NOT NULL,
  code varchar NOT NULL,
  display_order int,
  name varchar,
  hierarchy_id varchar NOT NULL,
  hierarchy_level_type_id varchar,
  parent varchar,
  FOREIGN KEY (parent) REFERENCES "hierarchy_entry"(id),
  FOREIGN KEY (hierarchy_id) REFERENCES "hierarchy"(id),
  FOREIGN KEY (hierarchy_level_type_id) REFERENCES "hierarchy_level_type"(id)
);

CREATE UNIQUE INDEX unq_hierarchy_entry_0 ON "hierarchy_entry"
(
  hierarchy_id,
  code
);


