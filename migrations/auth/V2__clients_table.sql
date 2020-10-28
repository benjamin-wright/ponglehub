BEGIN;

SAVEPOINT migration_2_restart;

DROP TABLE IF EXISTS clients;

CREATE TABLE clients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    callback_url VARCHAR(100) NOT NULL
);

RELEASE SAVEPOINT migration_2_restart;

COMMIT;