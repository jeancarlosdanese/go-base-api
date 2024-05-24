-- id	ptype	v0				v1					v2								v3	v4	v5
-- 1	p		admin			/api/v1/users		GET|POST						NULL	NULL	NULL
-- 2	p		admin			/api/v1/users/:id	GET|POST|PUT|PATCH|DELETE		NULL	NULL	NULL
CREATE OR REPLACE VIEW casbin_rules_view AS WITH policies AS (
        SELECT 'p' AS ptype,
            roles.name AS v0,
            endpoints.name AS v1,
            policies_roles.actions AS v2
        FROM policies_roles
            INNER JOIN roles ON policies_roles.role_id = roles.id
            INNER JOIN endpoints ON policies_roles.endpoint_id = endpoints.id
        UNION
        SELECT 'p' AS ptype,
            users.id::VARCHAR(36) AS v0,
            endpoints.name AS v1,
            policies_users.actions AS v2
        FROM policies_users
            INNER JOIN users ON policies_users.user_id = users.id
            INNER JOIN endpoints ON policies_users.endpoint_id = endpoints.id
    )
SELECT ROW_NUMBER() OVER (
        ORDER BY v0,
            v1
    ) AS id,
    ptype,
    v0,
    v1,
    v2,
    NULL AS v3,
    NULL AS v4,
    NULL AS v5
FROM policies
ORDER BY v0,
    v1;