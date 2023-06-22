COPY thi_arbre_validation (id, validation, sort_order, is_active) FROM '/Data/Go/sanarbo/cmd/sanarboServer/db/csv/thi_arbre_validation.csv' WITH DELIMITER ',' CSV HEADER;

COPY thi_arbre_to_be_checked (id, to_be_checked, sort_order, is_active) FROM '/Data/Go/sanarbo/cmd/sanarboServer/db/csv/thi_arbre_to_be_checked.csv' WITH DELIMITER ',' CSV HEADER;

COPY thi_arbre_note (id, note, sort_order, is_active) FROM '/Data/Go/sanarbo/cmd/sanarboServer/db/csv/thi_arbre_note.csv' WITH DELIMITER ',' CSV HEADER;

COPY thi_arbre_entourage (id, entourage, sort_order, is_active) FROM '/Data/Go/sanarbo/cmd/sanarboServer/db/csv/thi_arbre_entourage.csv' WITH DELIMITER ',' CSV HEADER;

COPY thi_arbre_chk (id, status, sort_order, is_active) FROM '/Data/Go/sanarbo/cmd/sanarboServer/db/csv/thi_arbre_chk.csv' WITH DELIMITER ',' CSV HEADER;

COPY thi_arbre_rev_surface (id, rev_surface, sort_order, is_active) FROM '/Data/Go/sanarbo/cmd/sanarboServer/db/csv/thi_arbre_rev_surface.csv' WITH DELIMITER ',' CSV HEADER;

COPY thi_arbre_etat_sanitaire (id, etat, sort_order, is_active) FROM '/Data/Go/sanarbo/cmd/sanarboServer/db/csv/thi_arbre_etat_sanitaire.csv' WITH DELIMITER ',' CSV HEADER;

COPY thi_arbre_etat_sanitaire_remarque (id, remarque, sort_order, is_active) FROM '/Data/Go/sanarbo/cmd/sanarboServer/db/csv/thi_arbre_etat_sanitaire_remarque.csv' WITH DELIMITER ',' CSV HEADER;
