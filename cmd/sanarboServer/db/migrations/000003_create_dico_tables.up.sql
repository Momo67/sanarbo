CREATE TABLE IF NOT EXISTS thi_arbre_validation
(
  id                      serial            constraint thi_arbre_validation_pk   primary key,
	validation              varchar(50)       not null,
	sort_order              int               not null,
  is_active               boolean           default true not null,
);

ALTER TABLE thi_arbre_validation OWNER TO sanarbo;
	
COMMENT ON TABLE thi_arbre_validation is 'thi_arbre_validation is a dico table containing values representing validation';

COPY thi_arbre_validation FROM '../csv/thi_arbre_validation.csv' WITH (FORMAT csv);

---
CREATE TABLE IF NOT EXISTS thi_arbre_to_be_checked
(
  id                      serial            constraint thi_arbre_to_be_checked_pk   primary key,
  to_be_checked           varchar(50)       not null,
  is_active               boolean           default true not null,
);

ALTER TABLE thi_arbre_to_be_checked OWNER TO sanarbo;
	
COMMENT ON TABLE thi_arbre_to_be_checked is 'thi_arbre_to_be_checked is a dico table containing values representing checking values';

---
CREATE TABLE IF NOT EXISTS thi_arbre_note
(
  id                      serial            constraint thi_arbre_note_pk   primary key,
  note                    int               not null,
	sort_order              int               not null,
  is_active               boolean           default true not null,
);

ALTER TABLE thi_arbre_note OWNER TO sanarbo;
	
COMMENT ON TABLE thi_arbre_note is 'thi_arbre_note is a dico table containing values representing notes';

---
CREATE TABLE IF NOT EXISTS thi_arbre_entourage
(
  id                      serial            constraint thi_arbre_entourage_pk   primary key,
	entourage               varchar(50)       not null,
	sort_order              int               not null,
  is_active               boolean           default true not null,
);

ALTER TABLE thi_arbre_entourage OWNER TO sanarbo;
	
COMMENT ON TABLE thi_arbre_entourage is 'thi_arbre_entourage is a dico table containing values representing entourage values';

---
CREATE TABLE IF NOT EXISTS thi_arbre_chk
(
  id                      serial            constraint thi_arbre_chk_pk   primary key,
	status                  varchar(10)       not null,
	sort_order              int               not null,
  is_active               boolean           default true not null,
);

ALTER TABLE thi_arbre_chk OWNER TO sanarbo;
	
COMMENT ON TABLE thi_arbre_chk is 'thi_arbre_chk is a dico table containing values representing checking values';

---
CREATE TABLE IF NOT EXISTS thi_arbre_rev_surface
(
  id                      serial            constraint thi_arbre_rev_surface_pk   primary key,
	rev_surface             varchar(50)       not null,
	sort_order              int               not null,
  is_active               boolean           default true not null,
);

ALTER TABLE thi_arbre_rev_surface OWNER TO sanarbo;
	
COMMENT ON TABLE thi_arbre_rev_surface is 'thi_arbre_rev_surface is a dico table containing values representing revetement surface values';

---
CREATE TABLE IF NOT EXISTS thi_arbre_etat_sanitaire
(
  id                      serial            constraint thi_arbre_etat_sanitaire_pk   primary key,
	etat                    varchar(30)       not null,
	sort_order              int               not null,
  is_active               boolean           default true not null,
);

ALTER TABLE thi_arbre_etat_sanitaire OWNER TO sanarbo;
	
COMMENT ON TABLE thi_arbre_etat_sanitaire is 'thi_arbre_etat_sanitaire is a dico table containing values representing état sanitaire values';

---
CREATE TABLE IF NOT EXISTS thi_arbre_etat_sanitaire_remarque
(
  id                      serial            constraint thi_arbre_etat_sanitaire_remarque_pk   primary key,
	remarque                varchar(100)      not null,
	sort_order              int               not null,
  is_active               boolean           default true not null,
);

ALTER TABLE thi_arbre_etat_sanitaire_remarque OWNER TO sanarbo;
	
COMMENT ON TABLE thi_arbre_etat_sanitaire_remarque is 'thi_arbre_etat_sanitaire_remarque is a dico table containing values representing remarques sur l''état sanitaire values';
