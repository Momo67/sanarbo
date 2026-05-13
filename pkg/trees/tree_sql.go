package trees

const (
	treesList = `
	SELECT id, name, description, is_active, create_time, creator, external_id, is_validated, ST_AsText(geom) as geom, json_build_object('idvalidation', tree_attributes::json->'idvalidation', 'ispublic', tree_attributes::json->'ispublic', 'essence', tree_attributes::json->'essence') as tree_att_light
	FROM tree_mobile
	WHERE is_active = true
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

	treesToValidate = `
	SELECT 
		t.id, 
		t.name, 
		t.description, 
		t.external_id, 
		t.is_validated, 
		t.last_modification_time, 
		COALESCE(u.name, '') AS last_modification_user, 
		ST_AsText(t.geom) as geom, 
		json_build_object(
			'idvalidation', t.tree_attributes::json->'idvalidation', 
			'ispublic', t.tree_attributes::json->'ispublic'
		) as tree_att_light
	FROM tree_mobile t
	LEFT JOIN go_user u ON u.external_id = t.last_modification_user
	LEFT JOIN lien_thing_thing arbre ON arbre.idthing1 = t.external_id
	LEFT JOIN lien_thing_thing sect ON sect.idthing1 = arbre.idthing2
	LEFT JOIN lien_thing_thing lien_empl ON lien_empl.idthing1 = t.external_id
	LEFT JOIN thing thing_empl ON thing_empl.idthing = lien_empl.idthing2 AND thing_empl.idtypething = 95
	LEFT JOIN lien_thing_thing sect_empl ON sect_empl.idthing1 = thing_empl.idthing
	WHERE t.is_validated = FALSE
	AND (
		($1::INTEGER IS NULL AND $2::INTEGER IS NULL)
		OR 
		($1::INTEGER IS NOT NULL AND $2::INTEGER IS NULL AND sect.idthing2 = $1::INTEGER)
		OR 
		($1::INTEGER IS NOT NULL AND $2::INTEGER IS NOT NULL 
		AND sect_empl.idthing2 = $1::INTEGER 
		AND lien_empl.idthing2 = $2::INTEGER)
		OR
		($1::INTEGER IS NULL AND $2::INTEGER IS NOT NULL AND lien_empl.idthing2 = $2::INTEGER)
	)
	GROUP BY 
		t.id, t.name, t.description, t.external_id, t.is_validated, 
		t.last_modification_time, u.name, t.geom, t.tree_attributes
	ORDER BY t.external_id;`

	validateTrees = `
	UPDATE tree_mobile
	SET is_validated = $2, id_validator = $3, datevalidation = NULL
	WHERE external_id = $1;`

	treesValidatedToUpdate = `
	SELECT xmlelement(name "Arbres", 
		xmlagg(
			xmlelement(name "ThiArbre",
				xmlforest(
					tm.external_id AS "IdObjet",
					tm.name AS "Nom",
					tm.description AS "Commentaire",
					74 AS "IdTypeThing",
					'arbre' AS "TypeThing",
					tm.last_modification_user AS "IdModificator",
					CASE WHEN tm.is_validated THEN 1 ELSE 0 END AS "IsValidated",
					NOW() AS "DateValidation",
					tm.tree_attributes->>'circonference' AS "Circonference",
					tm.tree_attributes->>'identourage' AS "IdEntourage",
					tm.tree_attributes->>'idchkentourage' AS "IdChkEntourage",
					tm.tree_attributes->>'entouragerem' AS "EntourageRem",
					tm.tree_attributes->>'idrevsurface' AS "IdRevSurface",
					tm.tree_attributes->>'idchkrevsurface' AS "IdChkRevSurface",
					tm.tree_attributes->>'revsurfacerem' AS "RevSurfaceRem",
					tm.tree_attributes->>'idetatsanitairepied' AS "IdEtatSanitairePied",
					tm.tree_attributes->>'idetatsanitairetronc' AS "IdEtatSanitaireTronc",
					tm.tree_attributes->>'idetatsanitairecouronne' AS "IdEtatSanitaireCouronne",
					tm.tree_attributes->>'idtobechecked' AS "IdToBeChecked",
					tm.tree_attributes->>'idvalidation' AS "IdValidation",
					tm.tree_attributes->>'idnote' AS "IdNote",
					tm.tree_attributes->>'etatsanitairerem' AS "EtatSanitaireRem",
				tm.tree_attributes->>'envracinairerem' AS "EnvRacinaireRem"
				)
			)
	)
	) AS xml_result
	FROM tree_mobile tm
	WHERE tm.is_validated = TRUE AND tm.datevalidation IS NULL;`

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
	  datevalidation		  timestamp without time zone,
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

	treesInsertFromGoeland = `INSERT INTO tree_mobile (name, description, external_id, is_active, inactivation_time, inactivation_reason, comment, is_validated, id_validator, datevalidation, create_time, creator, last_modification_time, last_modification_user, geom, tree_attributes)
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
                thing.datevalidation,
                thing.datecreated,
                thing.idcreator,
                COALESCE(thing.datelastmodif, '1970-01-01'),
                COALESCE(thing.idmodificator, 0),
                ST_GeomFromText(CONCAT('POINT(', to_char((thing_position.mineo/100.00), 'FM9999999.99'), ' ', to_char((thing_position.minsn/100.00), 'FM9999999.99'), ')'), 2056),
                (SELECT row_to_json(f) FROM (SELECT
                                                attr.idthing,
                                                attr.idvalidation,
                                                attr.ispublic,
                                                attr.ispublic,
                                                CASE
                                                    WHEN attr.idgenre = 127 OR attr.idespece IS NULL THEN '?'
                                                    ELSE UPPER(SUBSTRING(thi_arbre_genre.genre, 1, 1)) ||
                                                        SUBSTRING(
                                                            CASE
                                                                WHEN SUBSTRING(thi_arbre_espece.espece, 1, 1) = '''' THEN thi_arbre_espece.espece
                                                                ELSE thi_arbre_espece.espece
                                                            END,
                                                            CASE
                                                                WHEN SUBSTRING(thi_arbre_espece.espece, 1, 1) = '''' THEN 2
                                                                ELSE 1
                                                            END,
                                                            1
                                                        )
                                                END AS essence,
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
        LEFT JOIN thi_arbre_genre ON thi_arbre_genre.id = attr.idgenre
        LEFT JOIN thi_arbre_espece ON thi_arbre_espece.id = attr.idespece
        WHERE thi_arbre.idvalidation IN (1,5,6,7,8,9,10,11)
        ORDER BY thi_arbre.idthing
        ON CONFLICT (external_id) DO UPDATE
        SET name = EXCLUDED.name,
            description = EXCLUDED.description,
            datevalidation = EXCLUDED.datevalidation,
            tree_attributes = jsonb_set(tree_mobile.tree_attributes, '{idvalidation}', (EXCLUDED.tree_attributes->>'idvalidation')::jsonb);`

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
	WHERE thi_arbre.idthing = external_id AND is_validated = TRUE;`

	thingUpdate = `
	UPDATE thing
	SET
		idmodificator = tm.last_modification_user,
		datelastmodif = tm.last_modification_time,
		isvalidated = TRUE,
		datevalidation = NOW()
	FROM tree_mobile tm
	WHERE thing.idtypething = 74 AND thing.idthing = tm.external_id AND tm.is_validated = TRUE;`

	secteursList = `
	SELECT idthing as id, substring(name, 30, length(name)-29) as value
	FROM thing
	WHERE idtypething = 47 AND name LIKE '%Secteur d''entretien SPADOM%'
	ORDER BY value;`

	emplacementsList = `
	SELECT thing.idthing as id, substring(thing.name, 22, length(thing.name)-21) as value
	FROM thing
	INNER JOIN lien_thing_thing ON lien_thing_thing.idthing1 = thing.idthing
	WHERE thing.idtypething = 95
	ORDER BY value;`

	emplacementsListBySecteur = `
	SELECT thing.idthing as id, substring(thing.name, 22, length(thing.name)-21) as value
	FROM thing
	INNER JOIN lien_thing_thing ON lien_thing_thing.idthing1 = thing.idthing
	WHERE thing.idtypething = 95 AND lien_thing_thing.idthing2 = $1
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

	groupByName = `SELECT id, name, create_time, creator, last_modification_time, last_modification_user, is_active, inactivation_time, inactivation_reason, comment	
	FROM go_group
	WHERE name = $1;`
)
