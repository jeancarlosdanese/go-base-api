-- Entries Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS entries_id_seq;
-- Entries Table Definition
CREATE TABLE "public"."entries" (
    "id" int8 NOT NULL DEFAULT nextval('entries_id_seq'::regclass),
    "name" varchar(36) NOT NULL,
    PRIMARY KEY ("id")
);
-- Indices
CREATE UNIQUE INDEX uni_entries_name ON public.entries USING btree (name);
-- FIM Entries Table Definition
-- Permissions Sequence and defined types
CREATE SEQUENCE IF NOT EXISTS permissions_id_seq;
DROP TYPE IF EXISTS "public"."action_type";
CREATE TYPE "public"."action_type" AS ENUM ('list', 'store', 'show', 'update', 'delete');
-- Permissions Table Definition
CREATE TABLE "public"."permissions" (
    "id" int8 NOT NULL DEFAULT nextval('permissions_id_seq'::regclass),
    "entry_id" int8 NOT NULL,
    "action" "public"."action_type" NOT NULL,
    CONSTRAINT "fk_permissions_entries" FOREIGN KEY ("entry_id") REFERENCES "public"."entries"("id") ON DELETE RESTRICT ON UPDATE RESTRICT,
    PRIMARY KEY ("id")
);
-- Indices
CREATE UNIQUE INDEX uni_entries_entry_id_action ON public.permissions USING btree (entry_id, action);
-- FIM Permissions Table Definition
-- PermissionsRoles Table Definition
CREATE TABLE "public"."permissions_roles" (
    "permission_id" int8 NOT NULL,
    "role_id" int8 NOT NULL,
    CONSTRAINT "fk_permissions_roles_role" FOREIGN KEY ("role_id") REFERENCES "public"."roles"("id") ON DELETE RESTRICT ON UPDATE RESTRICT,
    CONSTRAINT "fk_permissions_roles_permission" FOREIGN KEY ("permission_id") REFERENCES "public"."permissions"("id") ON DELETE RESTRICT ON UPDATE RESTRICT,
    PRIMARY KEY ("permission_id", "role_id")
);
-- FIM PermissionsRoles Table Definition
-- PermissionsUsers Table Definition
CREATE TABLE "public"."permissions_users" (
    "permission_id" int8 NOT NULL,
    "user_id" uuid NOT NULL,
    CONSTRAINT "fk_permissions_users_user" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE RESTRICT ON UPDATE RESTRICT,
    CONSTRAINT "fk_permissions_users_permission" FOREIGN KEY ("permission_id") REFERENCES "public"."permissions"("id") ON DELETE RESTRICT ON UPDATE RESTRICT,
    PRIMARY KEY ("permission_id", "user_id")
);
-- FIM PermissionsUsers Table Definition