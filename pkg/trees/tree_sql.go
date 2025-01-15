package trees

const (
	treesList = `
	SELECT id, name, description, is_active, create_time, creator, external_id, is_validated, ST_AsText(geom) as geom, json_build_object('idvalidation', tree_attributes::json->'idvalidation', 'ispublic', tree_attributes::json->'ispublic') as tree_att_light
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
	  name                    text  not null 	CONSTRAINT name_min_length check (length(btrim(name)) > 2),
	  description             text           	CONSTRAINT description_min_length check (length(btrim(description)) > 0),
	  external_id             int 				CONSTRAINT unique_external_id UNIQUE (external_id),
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
	ALTER TABLE tree_mobile OWNER TO sanarbo;
	COMMENT ON TABLE tree_mobile is 'tree_mobile is the main table of the sanarbo application';`

	treesDicoGetValidation = "SELECT id, validation as value FROM thi_arbre_validation WHERE is_active = TRUE ORDER BY sort_order;"

	treesDicoGetToBeChecked = "SELECT id, to_be_checked as value FROM thi_arbre_to_be_checked WHERE is_active = TRUE ORDER BY sort_order;"

	treesDicoGetNote = "SELECT id, note::varchar(2) as value FROM thi_arbre_note WHERE is_active = TRUE ORDER BY sort_order;"

	treesDicoGetEntourage = "SELECT id, entourage as value FROM thi_arbre_entourage WHERE is_active = TRUE ORDER BY sort_order;"

	treesDicoGetChk = "SELECT id, status as value FROM thi_arbre_chk WHERE is_active = TRUE ORDER BY sort_order;"

	treesDicoGetRevSurface = "SELECT id, rev_surface as value FROM thi_arbre_rev_surface WHERE is_active = TRUE ORDER BY sort_order;"

	treesDicoGetEtatSanitaire = "SELECT id, etat as value FROM thi_arbre_etat_sanitaire WHERE is_active = TRUE ORDER BY sort_order;"
	
	treesDicoGetEtatSanitaireRem = "SELECT id, remarque as value FROM thi_arbre_etat_sanitaire_remarque WHERE is_active = TRUE ORDER BY sort_order;"

	treesInsertFromGoeland = `INSERT INTO tree_mobile (name, description, external_id, is_active, inactivation_time, inactivation_reason, comment, is_validated, id_validator, create_time, creator, last_modification_time, last_modification_user, geom, tree_attributes)
	SELECT
		REPLACE(thing.name, '''', ''''''),
		COALESCE(REPLACE(thing.description, '''', ''''''), NULL),
		thing.idthing,
		't',
		NULL,
		NULL,
		NULL,
		NULL,
		NULL,
		thing.datecreated,
		thing.idcreator,
		COALESCE(thing.datelastmodif, '1970-01-01'),
		COALESCE(thing.idmodificator, 0),
		ST_GeomFromText(CONCAT('POINT(', to_char((thing_position.mineo/100.00), 'FM9999999.99'), ' ', to_char((thing_position.minsn/100.00), 'FM9999999.99'), ')'), 2056),
		(SELECT row_to_json(f) FROM (SELECT 
						attr.idthing,
						attr.idvalidation,
						attr.ispublic,
						attr.idtobechecked,
						attr.idnote,
						attr.circonference,
						attr.identourage,
						attr.idchkentourage,
						attr.entouragerem,
						attr.idrevsurface,
						attr.idchkrevsurface,
						attr.revsurfacerem,
						attr.idetatsanitairepied,
						attr.idetatsanitairetronc,
						attr.idetatsanitairecouronne,
						attr.etatsanitairerem,
						attr.envracinairerem) f)
	FROM thi_arbre
	INNER JOIN thing ON thing.idthing = thi_arbre.idthing AND thing.isactive = true
	INNER JOIN thing_position ON thing_position.idthing = thi_arbre.idthing
	INNER JOIN thi_arbre attr ON attr.idthing = thing.idthing
	WHERE thi_arbre.idvalidation IN (1,5,6,7,8,9,10,11)
	ORDER BY thi_arbre.idthing
	ON CONFLICT (external_id) DO NOTHING;`

	thiArbreUpdate = `
	UPDATE thi_arbre
	SET 
		idthing = (tree_attributes->>'idthing')::INT,
		idvalidation = (tree_attributes->>'idvalidation')::INT,
		ispublic = (tree_attributes->>'ispublic')::BOOLEAN,
		idtobechecked = (tree_attributes->>'idtobechecked')::INT,
		idnote = (tree_attributes->>'idnote')::INT,
		circonference = (tree_attributes->>'circonference')::INT,
		identourage = (tree_attributes->>'identourage')::INT,
		idchkentourage = (tree_attributes->>'idchkentourage')::INT,
		entouragerem = tree_attributes->>'entouragerem',
		idrevsurface = (tree_attributes->>'idrevsurface')::INT,
		idchkrevsurface = (tree_attributes->>'idchkrevsurface')::INT,
		revsurfacerem = tree_attributes->>'revsurfacerem',
		idetatsanitairepied = (tree_attributes->>'idetatsanitairepied')::INT,
		idetatsanitairetronc = (tree_attributes->>'idetatsanitairetronc')::INT,
		idetatsanitairecouronne = (tree_attributes->>'idetatsanitairecouronne')::INT,
		etatsanitairerem = tree_attributes->>'etatsanitairerem',
		envracinairerem = tree_attributes->>'envracinairerem'
	FROM tree_mobile 
	WHERE thi_arbre.idthing = external_id AND is_validated = FALSE;`

	thingUpdate = `
	UPDATE thing
	SET
		idmodificator = tm.last_modification_user,
		datelastmodif = tm.last_modification_time,
		isvalidated = FALSE,
		datevalidation = NULL
	FROM tree_mobile tm
	WHERE thing.idtypething = 74 AND thing.idthing = tm.external_id AND tm.is_validated = FALSE;`

	secteursList = `WITH secteurs AS (
		SELECT DISTINCT UPPER(nom_sect) AS nom
		FROM geodata_gestion_com.spadom_surfaces
		ORDER BY nom
	)
	SELECT row_number() over () AS id, nom as value
	FROM secteurs;`

	emplacementsList = `SELECT DISTINCT idgo_empl AS id, SUBSTRING(nomgo_empl, LENGTH('Emplacement SPADOM - ') + 1) AS value
	FROM geodata_gestion_com.spadom_surfaces
	ORDER BY value;`

	emplacementsListBySecteur = `SELECT DISTINCT idgo_empl AS id, SUBSTRING(nomgo_empl, LENGTH('Emplacement SPADOM - ') + 1) AS value
	FROM geodata_gestion_com.spadom_surfaces
	WHERE UPPER(nom_sect) = $1
	ORDER BY value;`

	emplacementCentroid = `SELECT ST_ASText(ST_Centroid(ST_Collect(surface.the_geom))) AS geometry, ST_Area(ST_Collect(surface.the_geom)) AS surface
	FROM geodata_gestion_com.spadom_surfaces AS surface
	WHERE surface.idgo_empl = $1;`

	streetsList = `SELECT str.idthing AS id, str.lastname AS value, str.longname AS subtitle
	FROM thi_street str
	WHERE str.idville = 632
	ORDER BY str.lastname;`

	buildingsNumberByStreet = `SELECT sba.idaddress AS id, sba.number::text || COALESCE(sba.extention, '') AS value
	FROM thi_street_building_address sba
	WHERE sba.idthingstreet = $1;`

	buildingCenter = `SELECT 'POINT(' || (sba.coordeo / 100.0)::text || ' ' || (sba.coordsn / 100.0)::text || ')' AS geometry
	FROM thi_street_building_address sba
	WHERE sba.idaddress = $1`
)