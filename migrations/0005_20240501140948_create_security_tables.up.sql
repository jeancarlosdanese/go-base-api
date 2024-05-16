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