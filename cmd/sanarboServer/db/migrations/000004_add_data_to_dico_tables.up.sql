INSERT INTO thi_arbre_validation (id, validation, sort_order, is_active) VALUES 
  (1,'Existant',10,'t'),
  (2,'Supprimé',20,'t'),
  (3,'A repositionner',0,'f'),
  (4,'A replanter',0,'f'),
  (5,'En attente de soins',30,'t'),
  (6,'En attente d''abattage',40,'t'),
  (7,'En attente de remplacement',50,'t'),
  (8,'En attente de tomographie',60,'t'),
  (9,'A surveiller',70,'t'),
  (10,'En demande d''abattage',39,'t'),
  (11,'En attente de projet',55,'t');

INSERT INTO thi_arbre_to_be_checked (id, to_be_checked, sort_order, is_active) VALUES
  (1,'pour des raisons sanitaires',30,'t'),
  (2,'car inexistant',20,'t'),
  (3,'pour repositionnement',40,'t'),
  (4,'Non',10,'t');

INSERT INTO thi_arbre_note (id, note, sort_order, is_active) VALUES
  (1,1,10,'t'),
  (2,2,20,'t'),
  (3,3,30,'t'),
  (4,4,40,'t'),
  (5,5,50,'t');

INSERT INTO thi_arbre_entourage (id, entourage, sort_order, is_active) VALUES
  (1,'Bordures',10,'t'),
  (2,'Pavés',20,'t'),
  (3,'Métal',30,'t'),
  (4,'Autre',40,'t');

INSERT INTO thi_arbre_chk (id, status, sort_order, is_active) VALUES
  (1,'A vérifier',10,'t'),
  (2,'En ordre',20,'t'),
  (3,'A corriger',30,'t');

INSERT INTO thi_arbre_rev_surface (id, rev_surface, sort_order, is_active) VALUES
  (1,'Grille',10,'t'),
  (2,'Pavage',20,'t'),
  (3,'Minéral',30,'t'),
  (4,'Végétalisé',40,'t'),
  (5,'Autre',50,'t');

INSERT INTO thi_arbre_etat_sanitaire (id, etat, sort_order, is_active) VALUES
  (1,'En ordre',10,'t'),
  (2,'Blessé',20,'t'),
  (3,'Malade',30,'t');

INSERT INTO thi_arbre_etat_sanitaire_remarque (id, remarque, sort_order, is_active) VALUES
  (1,'Abattage',10,'t'),
  (2,'Haubanage',20,'t'),
  (3,'Taille',30,'t');
