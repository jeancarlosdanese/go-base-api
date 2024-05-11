-- Table Definition
CREATE TABLE "public"."users_roles" (
    "user_id" uuid NOT NULL,
    "role_id" int8 NOT NULL,
    CONSTRAINT "fk_users_roles_user" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE RESTRICT ON UPDATE RESTRICT,
    CONSTRAINT "fk_users_roles_role" FOREIGN KEY ("role_id") REFERENCES "public"."roles"("id") ON DELETE RESTRICT ON UPDATE RESTRICT,
    PRIMARY KEY ("user_id", "role_id")
);