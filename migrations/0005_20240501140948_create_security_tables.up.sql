-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS endpoints_id_seq;
-- Table Definition
CREATE TABLE "public"."endpoints" (
    "id" int8 NOT NULL DEFAULT nextval('endpoints_id_seq'::regclass),
    "name" varchar(254) NOT NULL,
    PRIMARY KEY ("id")
);
-- Indices
CREATE UNIQUE INDEX uni_endpoints_name ON public.endpoints USING btree (name);
-- FIM Endpoints Table Definition
-- Table Definition
CREATE TABLE "public"."policies_roles" (
    "role_id" int8 NOT NULL,
    "endpoint_id" int8 NOT NULL,
    "actions" varchar(36) NOT NULL,
    CONSTRAINT "fk_policies_roles_role" FOREIGN KEY ("role_id") REFERENCES "public"."roles"("id") ON DELETE RESTRICT ON UPDATE RESTRICT,
    CONSTRAINT "fk_policies_roles_endpoint" FOREIGN KEY ("endpoint_id") REFERENCES "public"."endpoints"("id") ON DELETE RESTRICT ON UPDATE RESTRICT,
    PRIMARY KEY ("role_id", "endpoint_id")
);
-- FIM PoliciesRoles Table Definition
-- Table Definition
CREATE TABLE "public"."policies_users" (
    "user_id" uuid NOT NULL,
    "endpoint_id" int8 NOT NULL,
    "actions" varchar(36) NOT NULL,
    CONSTRAINT "fk_policies_users_endpoint" FOREIGN KEY ("endpoint_id") REFERENCES "public"."endpoints"("id") ON DELETE RESTRICT ON UPDATE RESTRICT,
    CONSTRAINT "fk_policies_users_user" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE RESTRICT ON UPDATE RESTRICT,
    PRIMARY KEY ("user_id", "endpoint_id")
);
-- FIM PoliciesUsers Table Definition
-- Adicionar índices de performance
CREATE INDEX idx_policies_roles_role_id ON public.policies_roles USING btree (role_id);
CREATE INDEX idx_policies_roles_endpoint_id ON public.policies_roles USING btree (endpoint_id);
CREATE INDEX idx_policies_users_user_id ON public.policies_users USING btree (user_id);
CREATE INDEX idx_policies_users_endpoint_id ON public.policies_users USING btree (endpoint_id);
-- Inserir dados na tabela "endpoints" (só se não existir)
INSERT INTO "public"."endpoints" ("id", "name")
VALUES (1, '/api/v1/tenants'),
    (2, '/api/v1/tenants/:id'),
    (3, '/api/v1/users'),
    (4, '/api/v1/users/:id')
ON CONFLICT ("id") DO NOTHING;
-- Atualizar sequência da tabela "endpoints"
SELECT setval(
        'public.endpoints_id_seq',
        (
            SELECT MAX(id)
            FROM public.endpoints
        )
    );
-- Inserir dados na tabela "roles" (só se não existir)
INSERT INTO "public"."roles" ("id", "name")
VALUES (1, 'master'),
    (2, 'admin')
ON CONFLICT ("id") DO NOTHING;
-- Atualizar sequência da tabela "roles"
SELECT setval(
        'public.roles_id_seq',
        (
            SELECT MAX(id)
            FROM public.roles
        )
    );
-- Inserir dados na tabela "policies_roles" (só se não existir)
INSERT INTO "public"."policies_roles" ("role_id", "endpoint_id", "actions")
VALUES (1, 1, 'GET|POST'),
    (1, 2, 'GET|POST|PUT|PATCH|DELETE'),
    (1, 3, 'GET|POST'),
    (1, 4, 'GET|POST|PUT|PATCH|DELETE'),
    (2, 3, 'GET|POST'),
    (2, 4, 'GET|POST|PUT|PATCH|DELETE')
ON CONFLICT ("role_id", "endpoint_id") DO NOTHING;
-- Inserir dados do usuário "master" na tabela "users" (só se não existir)
-- NOTA: Hash bcrypt para senha 'master123' (cost=10)
-- Hash gerado dinamicamente: $2a$10$Y/QSjGygWZl49as2KLpJIeU5ENWeaM3.9D2GW7AAyzXJt.60jsVQm
-- Para gerar novo hash: htpasswd -nbBC 10 '' novasenha | cut -d: -f2
-- Ou usar Go: bcrypt.GenerateFromPassword([]byte("novasenha"), bcrypt.DefaultCost)
INSERT INTO "public"."users" (
        "tenant_id",
        "username",
        "name",
        "email",
        "password"
    )
SELECT
        t.id,
        'master',
        'Master User',
        'master@domain.local',
        '$2a$10$Y/QSjGygWZl49as2KLpJIeU5ENWeaM3.9D2GW7AAyzXJt.60jsVQm'
FROM tenants t
WHERE t.email = 'master@domain.local'
    AND NOT EXISTS (
        SELECT 1 FROM users u
        WHERE u.tenant_id = t.id
        AND u.email = 'master@domain.local'
    );
-- Atribuir role master para usuário master na tabela "users_roles" (só se não existir)
INSERT INTO "public"."users_roles" ("user_id", "role_id")
SELECT
        u.id,
        r.id
FROM users u
CROSS JOIN roles r
WHERE u.email = 'master@domain.local'
    AND r.name = 'master'
    AND NOT EXISTS (
        SELECT 1 FROM users_roles ur
        WHERE ur.user_id = u.id
        AND ur.role_id = r.id
    );

-- Conceder permissões de leitura de tenants para o role 'admin' (idempotente)
INSERT INTO "public"."policies_roles" ("role_id", "endpoint_id", "actions")
VALUES
    (2, 1, 'GET'),           -- admin pode listar tenants
    (2, 2, 'GET')            -- admin pode buscar tenant por ID
ON CONFLICT ("role_id", "endpoint_id") DO NOTHING;

-- Criar usuário 'admin' (idempotente) utilizando o mesmo tenant do usuário master
INSERT INTO "public"."users" (
        "tenant_id",
        "username",
        "name",
        "email",
        "password"
    )
SELECT
        t.id,
        'admin',
        'Admin User',
        'admin@domain.local',
        '$2a$10$Y/QSjGygWZl49as2KLpJIeU5ENWeaM3.9D2GW7AAyzXJt.60jsVQm'
FROM tenants t
WHERE t.email = 'master@domain.local'
    AND NOT EXISTS (
        SELECT 1 FROM users u
        WHERE u.tenant_id = t.id
        AND u.email = 'admin@domain.local'
    );

-- Atribuir role 'admin' para o usuário admin (idempotente)
INSERT INTO "public"."users_roles" ("user_id", "role_id")
SELECT
        u.id,
        r.id
FROM users u
CROSS JOIN roles r
WHERE u.email = 'admin@domain.local'
    AND r.name = 'admin'
    AND NOT EXISTS (
        SELECT 1 FROM users_roles ur
        WHERE ur.user_id = u.id
        AND ur.role_id = r.id
    );