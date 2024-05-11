-- Desfazer a criação da tabela 'roles' e remover o índice e a sequência associados
BEGIN;
-- Remove a tabelas se ela existirem
DROP TABLE IF EXISTS "public"."permissions_users";
DROP TABLE IF EXISTS "public"."permissions_roles";
DROP TABLE IF EXISTS "public"."permissions";
DROP TABLE IF EXISTS "public"."entries";
-- Remove a types criados para a geração de tables se eles existirem
DROP TYPE IF EXISTS "public"."action_type";
-- Remove a sequências criadas para a geração de IDs se ela existirem
DROP SEQUENCE IF EXISTS permissions_id_seq;
DROP SEQUENCE IF EXISTS entries_id_seq;
-- Os índices serão automaticamente removidos com as tabelas, mas se precisar removê-lo explicitamente:
-- DROP INDEX IF EXISTS public.uni_roles_name;
COMMIT;