create table vacancies
(
    id          serial
        constraint vacancies_pk
            primary key,
    profile_id  int       not null
        constraint vacancies_profiles_id_fk
            references profiles,
    title       varchar   not null,
    address varchar not null,
    lat         float     not null,
    long        float     not null,
    description text      not null,
    created_at  timestamp not null,
    updated_at  timestamp,
    deleted_at  timestamp
);
