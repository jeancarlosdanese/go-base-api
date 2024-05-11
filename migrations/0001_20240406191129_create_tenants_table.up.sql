-- This script only contains the table creation statements and does not fully represent the table in the database. Do not use it as a backup.
DROP TYPE IF EXISTS "public"."person_type";
CREATE TYPE "public"."person_type" AS ENUM ('FISICA', 'JURIDICA');
DROP TYPE IF EXISTS "public"."status_type";
CREATE TYPE "public"."status_type" AS ENUM ('ATIVO', 'INATIVO');
-- Table Definition
CREATE TABLE "public"."tenants" (
    "id" uuid NOT NULL DEFAULT gen_random_uuid(),
    "created_at" timestamptz,
    "updated_at" timestamptz,
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
    "subdomain" varchar(20) NOT NULL,
    "domain" varchar(100),
    "status" "public"."status_type" NOT NULL DEFAULT 'ATIVO'::status_type,
    PRIMARY KEY ("id")
);
-- Indices
CREATE UNIQUE INDEX uni_tenants_cpf_cnpj ON public.tenants USING btree (cpf_cnpj);
CREATE UNIQUE INDEX uni_tenants_subdomain_domain ON public.tenants USING btree (subdomain, domain);
CREATE INDEX idx_tenants_deleted_at ON public.tenants USING btree (deleted_at);