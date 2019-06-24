create table books (
    id serial primary key,
    title text not null,
    author text not null,
    year text
);
create table users (
    id serial primary key,
    email text not null unique,
    password text not null
);
