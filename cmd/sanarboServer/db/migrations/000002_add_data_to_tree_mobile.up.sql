INSERT INTO tree_mobile
(name, description, is_active, create_time, creator, geom, tree_attributes) 
VALUES ('MyNewTree', 'Test de création', 't', CURRENT_TIMESTAMP, 999, ST_GeomFromText('POINT(2538221 1152372)', 2056), '{"idvalidation":1, "etatsanitairerem":"Rien à signaler"}');

INSERT INTO tree_mobile
(name, description, is_active, create_time, creator, geom, tree_attributes) 
VALUES ('Mon bel arbre', 'création d''un nouvel arbre', 't', CURRENT_TIMESTAMP, 999, ST_GeomFromText('POINT(2538221 1152372)', 2056), '{"idvalidation":2, "enracinairerem":"En ordre"}');

INSERT INTO tree_mobile
  (
    name, 
    description, 
    external_id, 
    is_active, 
    inactivation_time, 
    inactivation_reason, 
    comment, 
    is_validated, 
    id_validator, 
    create_time, 
    creator, 
    last_modification_time, 
    last_modification_user, 
    geom, 
    tree_attributes
  )
VALUES 
  (
    '##TEST MPittet & CGil ## Wollemia nobilis - secteur OUEST - SI Chauderon 23-25-27 (70852)', 
    'Anciennement: Wollemia nobilis - secteur OUEST - SI Chauderon 23-25-27 (70852)',
    101357,
    't', 
    NULL, 
    NULL, 
    NULL, 
    NULL, 
    NULL, 
    '2013-10-28 09:39:00', 
    10957, 
    '2015-09-17 23:01:00', 
    10957, 
    ST_GeomFromText('POINT(2537600.37 1152644.94)', 2056),
    '{"idvalidation":1,
    "idtobechecked":4,
    "idnote":1,
    "circonference":666,
    "identourage":4,
    "idchkentourage":2,
    "entouragerem":"a soigner",
    "idrevsurface":5,
    "idchkrevsurface":2,
    "revsurfacerem":"",
    "idetatsanitairepied":2,
    "idetatsanitairetronc":1,
    "idetatsanitairecouronne":1,
    "etatsanitairerem":"Taille. Racines un peu rouillées... Ca mériterait un coup de pinceau...",
    "envracinairerem":""}'
    )
    