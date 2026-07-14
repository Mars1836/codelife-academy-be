-- The migration runner bootstraps schema_migrations before applying files.
-- This first migration marks the beginning of the application schema history.
SELECT 1;
