CREATE TABLE users (
    id serial primary key,
    email varchar(150) not null unique,
    password text not null
);