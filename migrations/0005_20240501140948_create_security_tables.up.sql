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
-- Inserir dados na tabela "endpoints"
INSERT INTO "public"."endpoints" ("id", "name")
VALUES (1, '/api/v1/tenants'),
    (2, '/api/v1/tenants/:id'),
    (3, '/api/v1/users'),
    (4, '/api/v1/users/:id');
-- Atualizar sequência da tabela "endpoints"
SELECT setval(
        'public.endpoints_id_seq',
        (
            SELECT MAX(id)
            FROM public.endpoints
        )
    );
-- Inserir dados na tabela "roles"
INSERT INTO "public"."roles" ("id", "name")
VALUES (1, 'master'),
    (2, 'admin');
-- Atualizar sequência da tabela "roles"
SELECT setval(
        'public.roles_id_seq',
        (
            SELECT MAX(id)
            FROM public.roles
        )
    );
-- Inserir dados na tabela "policies_roles"
INSERT INTO "public"."policies_roles" ("role_id", "endpoint_id", "actions")
VALUES (1, 1, 'GET|POST'),
    (1, 2, 'GET|POST|PUT|PATCH|DELETE'),
    (1, 3, 'GET|POST'),
    (1, 4, 'GET|POST|PUT|PATCH|DELETE'),
    (2, 3, 'GET|POST'),
    (2, 4, 'GET|POST|PUT|PATCH|DELETE');
-- Inserir dados do usuário "master" na tabela "users"
INSERT INTO "public"."users" (
        "tenant_id",
        "username",
        "name",
        "email",
        "password"
    )
VALUES (
        (
            SELECT id
            FROM tenants
            WHERE name = 'Master Tenant'
        ),
        'master',
        'Master',
        'master@domain.local',
        '$2a$10$h16kftLEZP/YYQtvv7xjQutiknBJuNHhnZWWlfSG6WFlwvtg5i3MG'
    );
-- Atribuir role master para usuário master na tabela "users_roles"
INSERT INTO "public"."users_roles" ("user_id", "role_id")
VALUES (
        (
            SELECT id
            FROM users
            WHERE email = 'master@domain.local'
        ),
        (
            SELECT id
            FROM roles
            WHERE name = 'master'
        )
    );