create table profiles
(
    id serial not null
        constraint profile_pk
            primary key,
    mobile varchar not null,
    email varchar(255),
    password varchar(255) not null,
    name varchar,
    middle_name varchar,
    sure_name varchar,
    is_active bool,
    roles json,
    created_at timestamp not null,
    updated_at timestamp,
    deleted_at timestamp
);

create unique index profile_email_uindex
    on profiles (email);