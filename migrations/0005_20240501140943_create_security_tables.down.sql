-- Desfazer a criação da tabela 'roles' e remover o índice e a sequência associados
BEGIN;
-- Remove a tabelas se ela existirem
DROP TABLE IF EXISTS "public"."policies_users";
DROP TABLE IF EXISTS "public"."policies_roles";
DROP TABLE IF EXISTS "public"."endpoints";
-- Remove a types criados para a geração de tables se eles existirem
DROP TYPE IF EXISTS "public"."action_type";
-- Remove a sequências criadas para a geração de IDs se ela existirem
DROP SEQUENCE IF EXISTS endpoints_id_seq;
-- Os índices serão automaticamente removidos com as tabelas, mas se precisar removê-lo explicitamente:
-- DROP INDEX IF EXISTS public.uni_roles_name;
COMMIT;