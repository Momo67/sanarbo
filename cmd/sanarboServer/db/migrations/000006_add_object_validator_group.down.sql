-- Only remove the group when nobody belongs to it, so a group created by hand and in use survives a rollback.
DELETE FROM public.go_group g
WHERE g.name = 'object_validator'
  AND NOT EXISTS (SELECT 1 FROM public.go_user u WHERE g.id = ANY (u.groups_id));
