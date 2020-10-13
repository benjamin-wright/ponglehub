BEGIN;

SAVEPOINT migration_1_restart;

DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password VARCHAR(100) NOT NULL,
    verified BOOLEAN NOT NULL
);

RELEASE SAVEPOINT migration_1_restart;

COMMIT;