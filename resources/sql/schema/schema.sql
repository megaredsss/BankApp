CREATE TABLE users (
    id serial primary key,
    email varchar(150) not null unique,
    password text not null
);

CREATE TABLE users_profiles(
    id serial primary key,
    users_id integer not null REFERENCES users(id) ON DELETE CASCADE,
    first_name varchar(20) not null,
    last_name varchar(20) not null,
    balance numeric(10, 2) not null default 0.00
);