package trees

const (
	treesList = `
	SELECT id, name, description, is_active, create_time, creator, external_id, ST_AsText(geom) as geom, json_build_object('idvalidation', tree_attributes::json->'idvalidation') as tree_att_light
	FROM tree_mobile
	LIMIT $1 OFFSET $2;`

	treesGet = `
	SELECT id, name, description, external_id, is_active, inactivation_time, inactivation_reason, comment, is_validated, id_validator,
			create_time, creator, last_modification_time, last_modification_user, ST_AsText(geom) as geom, tree_attributes
	FROM tree_mobile
	WHERE id = $1;`
	
	treesGetMaxId = "SELECT MAX(id) FROM tree_mobile;"

	treesExist = "SELECT COUNT(*) FROM tree_mobile WHERE id = $1;" 

	treesCount = "SELECT COUNT(*) FROM tree_mobile;"

	treesCreate = `
	INSERT INTO tree_mobile
	(name, description, external_id, is_active, comment, create_time, creator, geom, tree_attributes) 
	VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, $6, ST_GeomFromText($7, 2056), $8)
	RETURNING id;`
	
	treesUpdate = `
	UPDATE tree_mobile
	SET name					= $1,
		description				= $2,
		external_id				= $3,
		is_active				= $4,
		inactivation_time		= $5,
		inactivation_reason		= $6,
		comment					= $7,
		is_validated			= $8,
		id_validator			= $9,
		last_modification_time 	= CURRENT_TIMESTAMP,
		last_modification_user	= $10,
		geom					= ST_GeomFromText($11, 2056),
		tree_attributes			= $12
	WHERE id = $13;`

	treesDelete = "DELETE FROM tree_mobile WHERE id = $1;"

	treesSearchByName = "SELECT id, name, description, is_active, create_time, creator, external_id FROM tree_mobile WHERE name LIKE $1;"

	treesIsActive = "SELECT isactive FROM tree_mobile WHERE id = $1;"

	treesCreateTable = `
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
	  create_time             timestamp default now() not null,
	  creator                 integer  not null,
	  last_modification_time  timestamp,
	  last_modification_user  integer,
	  geom                    geometry(Point,2056)  not null,
	  tree_attributes         jsonb not null
	);
	ALTER TABLE tree_mobile OWNER TO postgres;
	COMMENT ON TABLE tree_mobile is 'tree_mobile is the main table of the sanarbo application';`

	treesDicoGetValidation = "SELECT id, validation as value FROM thi_arbre_validation WHERE is_active = TRUE ORDER BY sort_order;"

	treesDicoGetToBeChecked = "SELECT id, to_be_checked as value FROM thi_arbre_to_be_checked WHERE is_active = TRUE ORDER BY sort_order;"

	treesDicoGetNote = "SELECT id, note::varchar(2) as value FROM thi_arbre_note WHERE is_active = TRUE ORDER BY sort_order;"

	treesDicoGetEntourage = "SELECT id, entourage as value FROM thi_arbre_entourage WHERE is_active = TRUE ORDER BY sort_order;"

	treesDicoGetChk = "SELECT id, status as value FROM thi_arbre_chk WHERE is_active = TRUE ORDER BY sort_order;"

	treesDicoGetRevSurface = "SELECT id, rev_surface as value FROM thi_arbre_rev_surface WHERE is_active = TRUE ORDER BY sort_order;"

	treesDicoGetEtatSanitaire = "SELECT id, etat as value FROM thi_arbre_etat_sanitaire WHERE is_active = TRUE ORDER BY sort_order;"
	
	treesDicoGetEtatSanitaireRem = "SELECT id, remarque as value FROM thi_arbre_etat_sanitaire_remarque WHERE is_active = TRUE ORDER BY sort_order;"

	treesInsertFromGoeland = `SELECT thi_arbre.idvalidation, 'INSERT INTO tree_mobile (name, description, external_id, is_active, inactivation_time, inactivation_reason, comment, is_validated, id_validator, create_time, creator, last_modification_time, last_modification_user, geom, tree_attributes) 
	VALUES (''' 
		|| REPLACE(thing.name, '''', '''''')
		|| ''',' || COALESCE('''' || REPLACE(thing.description, '''', '''''') || '''', 'NULL') 
		|| ',' || thing.idthing
		|| ',''t'''
		|| ',NULL'
		|| ',NULL'
		|| ',NULL'
		|| ',' || COALESCE(thing.isvalidated,'f')
		|| ',NULL'
		|| ',''' || thing.datecreated || ''''
		|| ',' || thing.idcreator
		|| ',''' || thing.datelastmodif || ''''
		|| ',' || thing.idmodificator
		|| ',ST_GeomFromText(''POINT(2' || to_char((thing_position.mineo/100.00), 'FM9999999.99') || ' 1' || to_char((thing_position.minsn/100.00), 'FM9999999.99') || ')'', 2056)'
		|| ',''' || json_build_object('idvalidation', COALESCE(thi_arbre.idvalidation::text,'')::integer,
									  'idtobechecked', COALESCE(thi_arbre.idtobechecked::text,'')::integer,
									  'idnote', COALESCE(thi_arbre.idnote::text,'')::integer,
									  'circonference', COALESCE(thi_arbre.circonference::varchar(10),'')::integer,
									  'identourage', COALESCE(thi_arbre.identourage::text,'')::integer,
									  'idchkentourage', COALESCE(thi_arbre.idchkentourage::text,'')::integer,
									  'entouragerem', COALESCE(thi_arbre.entouragerem::text,''),
									  'idrevsurface', COALESCE(thi_arbre.idrevsurface::text,'')::integer,
									  'idchkrevsurface', COALESCE(thi_arbre.idchkrevsurface::text,'')::integer,
									  'revsurfacerem', COALESCE(thi_arbre.revsurfacerem::text,''),
									  'idetatsanitairepied', COALESCE(thi_arbre.idetatsanitairepied::text,'')::integer,
									  'idetatsanitairetronc', COALESCE(thi_arbre.idetatsanitairetronc::text,'')::integer,
									  'idetatsanitairecouronne', COALESCE(thi_arbre.idetatsanitairecouronne::text,'')::integer,
									  'etatsanitairerem', COALESCE(thi_arbre.etatsanitairerem::text,''),
									  'envracinairerem', COALESCE(thi_arbre.envracinairerem::text,''))
		|| ''');'
	FROM thi_arbre
	INNER JOIN thing ON thing.idthing = thi_arbre.idthing AND thing.isactive = true
	INNER JOIN thing_position ON thing_position.idthing = thi_arbre.idthing
	WHERE thi_arbre.idvalidation = 5
	LIMIT 1;`
)
