-- Desfazer a criação da tabela 'roles' e remover o índice e a sequência associados
BEGIN;
-- Remove a tabela se ela existir
DROP TABLE IF EXISTS public.roles;
-- Remove a sequência criada para a geração de IDs se ela existir
DROP SEQUENCE IF EXISTS roles_id_seq;
-- O índice será automaticamente removido com a tabela, mas se precisar removê-lo explicitamente:
-- DROP INDEX IF EXISTS public.uni_roles_name;
COMMIT;