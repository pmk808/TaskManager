CREATE SCHEMA IF NOT EXISTS task_management;

GRANT ALL PRIVILEGES ON SCHEMA task_management TO postgres;

SET search_path TO task_management;

ALTER DEFAULT PRIVILEGES IN SCHEMA task_management
    GRANT ALL ON TABLES TO postgres;

ALTER DEFAULT PRIVILEGES IN SCHEMA task_management
    GRANT ALL ON SEQUENCES TO postgres;