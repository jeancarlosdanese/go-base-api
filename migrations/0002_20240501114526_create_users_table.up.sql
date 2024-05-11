-- Table Definition
CREATE TABLE "public"."users" (
    "id" uuid NOT NULL DEFAULT gen_random_uuid(),
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "tenant_id" uuid NOT NULL,
    "username" varchar(80) NOT NULL,
    "name" varchar(254) NOT NULL,
    "email" varchar(100) NOT NULL,
    "password" varchar(60) NOT NULL,
    "thumbnail" varchar(70),
    CONSTRAINT "fk_users_tenant" FOREIGN KEY ("tenant_id") REFERENCES "public"."tenants"("id") ON DELETE RESTRICT ON UPDATE RESTRICT,
    PRIMARY KEY ("id")
);
-- Indices
CREATE UNIQUE INDEX uni_users_tenant_id_email ON public.users USING btree (tenant_id, email);
CREATE INDEX idx_users_deleted_at ON public.users USING btree (deleted_at);