-- Sequence and defined type
-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS roles_id_seq;
-- Table Definition
CREATE TABLE "public"."roles" (
    "id" int8 NOT NULL DEFAULT nextval('roles_id_seq'::regclass),
    "name" varchar(36) NOT NULL,
    PRIMARY KEY ("id")
);
-- Indices
CREATE UNIQUE INDEX uni_roles_name ON public.roles USING btree (name);