drop table if exists users;
drop table if exists used_access_tokens;
drop table if exists user_refresh_tokens;

create table users (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name varchar(100) not null,
    email varchar(100) not null,
    password varchar(100) not null,
    verified boolean not null,
);

create table used_access_tokens (
    id uuid not null PRIMARY KEY,
    user_id uuid not null,
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
            REFERENCES users(id)
            ON DELETE CASCADE
);

create table used_refresh_tokens (
    id uuid not null PRIMARY KEY,
    user_id uuid not null,
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
            REFERENCES users(id)
            ON DELETE CASCADE
);