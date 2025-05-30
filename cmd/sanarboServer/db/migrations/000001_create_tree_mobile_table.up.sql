CREATE TABLE IF NOT EXISTS tree_mobile
(
  id                      serial            CONSTRAINT tree_mobile_pk   primary key,
  name                    text  not null constraint name_min_length check (length(btrim(name)) > 2),
  description             text           constraint description_min_length check (length(btrim(description)) > 2),
  external_id             int,
  is_active               boolean default true not null,
  inactivation_time       timestamp,
  inactivation_reason     text,
  comment                 text,
  is_validated            boolean default false,
  id_validator            int,
  datevalidation    		  timestamp without time zone,
  create_time             timestamp default now() not null,
  creator                 integer  not null,
  last_modification_time  timestamp,
  last_modification_user  integer,
  geom                    geometry(Point,2056)  not null,
  tree_attributes         jsonb not null
);
CREATE INDEX tree_mobile_geom_idx ON tree_mobile USING GIST (geom);

COMMENT ON TABLE tree_mobile is 'tree_mobile is the main table of the sanarbo application';
