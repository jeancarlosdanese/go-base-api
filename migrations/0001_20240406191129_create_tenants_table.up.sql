DROP TYPE IF EXISTS "public"."person_type";
CREATE TYPE "public"."person_type" AS ENUM ('FISICA', 'JURIDICA');
DROP TYPE IF EXISTS "public"."status_type";
CREATE TYPE "public"."status_type" AS ENUM ('ATIVO', 'INATIVO');
-- Table Definition
CREATE TABLE "public"."tenants" (
    "id" uuid NOT NULL DEFAULT gen_random_uuid(),
    "created_at" timestamptz DEFAULT now(),
    "updated_at" timestamptz DEFAULT now(),
    "deleted_at" timestamptz,
    "type" "public"."person_type" NOT NULL,
    "name" varchar(100) NOT NULL,
    "cpf_cnpj" varchar(18),
    "ie" varchar(20),
    "cep" varchar(9),
    "street" varchar(100),
    "number" varchar(10),
    "neighborhood" varchar(100),
    "city" varchar(100),
    "state" varchar(2),
    "complement" varchar(100),
    "email" varchar(100),
    "phone" varchar(15),
    "cell_phone" varchar(15),
    "api_key" varchar(254),
    "allowed_origins" jsonb,
    "status" "public"."status_type" NOT NULL DEFAULT 'ATIVO'::status_type,
    PRIMARY KEY ("id")
);
-- Indices
CREATE UNIQUE INDEX uni_tenants_cpf_cnpj ON public.tenants USING btree (cpf_cnpj) WHERE cpf_cnpj IS NOT NULL;
CREATE UNIQUE INDEX uni_tenants_email ON public.tenants USING btree (email) WHERE email IS NOT NULL;
CREATE UNIQUE INDEX uni_tenants_api_key ON public.tenants USING btree (api_key) WHERE api_key IS NOT NULL;
CREATE INDEX idx_tenants_deleted_at ON public.tenants USING btree (deleted_at);
CREATE INDEX idx_tenants_status ON public.tenants USING btree (status);
-- Inserir dados do "master" Tenant na tabela "tenants"
INSERT INTO "public"."tenants" (
        "type",
        "name",
        "email",
        "api_key",
        "allowed_origins",
        "status"
    )
VALUES (
        'JURIDICA',
        'Master Tenant',
        'master@domain.local',
        gen_random_uuid()::text,
        '["localhost"]'::jsonb,
        'ATIVO'
    );