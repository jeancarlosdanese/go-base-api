-- id	ptype	v0				v1		v2		v3	v4	v5
-- 1	p		administrator	user	list
-- 2	p		administrator	user	store
-- 3	p		administrator	user	show
CREATE OR REPLACE VIEW casbin_rules_view AS
SELECT row_number() OVER (
        ORDER BY roles.name,
            entries.name,
            permissions.action
    ) AS id,
    'p' ptype,
    roles.name v0,
    entries.name v1,
    permissions.action v2,
    null v3,
    null v4,
    null v5
FROM permissions_roles
    INNER JOIN roles ON permissions_roles.role_id = roles.id
    INNER JOIN permissions ON permissions_roles.permission_id = permissions.id
    INNER JOIN entries ON permissions.entry_id = entries.id;
SELECT *
FROM casbin_rules_view;