-- One-off fix: apostrophes doubled in tree_mobile.name and tree_mobile.description by the Goeland import.
--
-- treesInsertFromGoeland (pkg/trees/tree_sql.go) used REPLACE(x, '''', '''''') in an INSERT ... SELECT,
-- so every apostrophe coming from Goeland (thing.name, thing.description) was stored twice: 'Euchlora' -> ''Euchlora''.
-- This script restores the Goeland value, only for rows equal to the doubled Goeland value:
-- names or descriptions edited in the mobile app are left untouched.
-- last_modification_time and last_modification_user are not changed: this is a data repair, not a user edit.
--
-- Fix the import query (and the job running it) first, otherwise the next import doubles the apostrophes again.
--
-- Dry run (default): applies the fix inside a transaction, shows the counts, then ROLLBACK.
--   psql -h <host> -U sanarbo -d sanarbo -f scripts/fix_doubled_apostrophes.sql
-- Apply: same command with -v apply=1, the transaction is committed.
--   psql -h <host> -U sanarbo -d sanarbo -v apply=1 -f scripts/fix_doubled_apostrophes.sql

\set ON_ERROR_STOP on

BEGIN;

\echo '== before'
SELECT count(*) FILTER (WHERE tm.name <> t.name
                          AND tm.name = replace(t.name, '''', '''''')) AS names_doubled,
       count(*) FILTER (WHERE tm.description IS DISTINCT FROM t.description
                          AND tm.description = replace(t.description, '''', '''''')) AS descriptions_doubled
FROM tree_mobile tm
JOIN thing t ON t.idthing = tm.external_id AND t.idtypething = 74;

\echo '== fixing names'
UPDATE tree_mobile tm
SET name = t.name
FROM thing t
WHERE t.idthing = tm.external_id
  AND t.idtypething = 74
  AND tm.name <> t.name
  AND tm.name = replace(t.name, '''', '''''');

\echo '== fixing descriptions'
UPDATE tree_mobile tm
SET description = t.description
FROM thing t
WHERE t.idthing = tm.external_id
  AND t.idtypething = 74
  AND tm.description IS DISTINCT FROM t.description
  AND tm.description = replace(t.description, '''', '''''');

\echo '== after (both counts must be 0)'
SELECT count(*) FILTER (WHERE tm.name <> t.name
                          AND tm.name = replace(t.name, '''', '''''')) AS names_doubled,
       count(*) FILTER (WHERE tm.description IS DISTINCT FROM t.description
                          AND tm.description = replace(t.description, '''', '''''')) AS descriptions_doubled
FROM tree_mobile tm
JOIN thing t ON t.idthing = tm.external_id AND t.idtypething = 74;

\if :{?apply}
  COMMIT;
  \echo '== COMMITTED'
\else
  ROLLBACK;
  \echo '== DRY RUN: rolled back, nothing was changed. Run again with -v apply=1 to apply.'
\endif
