SELECT pg_terminate_backend(pid) 
FROM pg_stat_activity 
WHERE datname = 'taskmanager'
AND pid <> pg_backend_pid();

DROP DATABASE IF EXISTS taskmanager;

CREATE DATABASE taskmanager
    WITH 
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'English_Philippines.1252'
    LC_CTYPE = 'English_Philippines.1252'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;