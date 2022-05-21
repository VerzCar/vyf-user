BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE gender AS ENUM (
    'X',
    'MAN',
    'WOMEN'
    );

create table countries
(
    id             bigserial
        constraint countries_pkey
            primary key,
    name           varchar(255) not null,
    alpha2         char(2)      not null,
    alpha3         char(3)      not null,
    continent_code char(2)      not null,
    number         char(3)      not null,
    full_name      varchar(255) not null,
    created_at     timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at     timestamp with time zone default CURRENT_TIMESTAMP
);

create index idx_countries_id
    on countries (id);

create table contacts
(
    id                       bigserial
        constraint contacts_pkey
            primary key,
    email                    varchar(150) not null
        constraint chk_contacts_email
            check ((email)::text <> ''::text),
    phone_number             varchar(20)  not null,
    phone_number_country_id  bigint       not null
        constraint fk_contacts_phone_number_country
            references countries
            on delete restrict,
    phone_number2            varchar(20),
    phone_number2_country_id bigint
        constraint fk_contacts_phone_number2_country
            references countries
            on delete restrict,
    web                      text,
    created_at               timestamp with time zone,
    updated_at               timestamp with time zone
);

create index idx_contacts_email
    on contacts (email);

create index idx_contacts_id
    on contacts (id);

create table addresses
(
    id          bigserial
        constraint addresses_pkey
            primary key,
    address     varchar(100) not null
        constraint chk_addresses_address
            check ((address)::text <> ''::text),
    city        varchar(80)  not null
        constraint chk_addresses_city
            check ((city)::text <> ''::text),
    postal_code varchar(15)  not null
        constraint chk_addresses_postal_code
            check ((postal_code)::text <> ''::text),
    country_id  bigint       not null
        constraint fk_addresses_country
            references countries
            on delete restrict,
    created_at  timestamp with time zone,
    updated_at  timestamp with time zone
);

create index idx_addresses_id
    on addresses (id);

create table locales
(
    id            bigserial
        constraint locales_pkey
            primary key,
    locale        varchar(29) not null,
    language_code varchar(9),
    lcid_string   varchar(9),
    created_at    timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at    timestamp with time zone default CURRENT_TIMESTAMP
);

create index idx_locales_id
    on locales (id);

create table companies
(
    id             bigserial
        constraint companies_pkey
            primary key,
    name           varchar(100)          not null
        constraint chk_companies_name
            check ((name)::text <> ''::text),
    address_id     bigint                not null
        constraint fk_companies_address
            references addresses
            on delete cascade,
    contact_id     bigint                not null
        constraint fk_companies_contact
            references contacts
            on delete cascade,
    owner_id       bigint,
    brand_logo_url text,
    tax_id         varchar(20),
    is_verified    boolean default false not null,
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone
);

create index idx_companies_id
    on companies (id);

create table users
(
    id          bigserial
        constraint users_pkey
            primary key,
    identity_id varchar(50) unique         not null,
    first_name  varchar(50),
    last_name   varchar(50),
    gender      gender default 'X'::gender not null,
    locale_id   bigint
        constraint fk_users_locale
            references locales
            on delete restrict,
    address_id  bigint
        constraint fk_users_address
            references addresses
            on delete cascade,
    contact_id  bigint
        constraint fk_users_contact
            references contacts
            on delete cascade,
    company_id  bigint
        constraint fk_companies_users
            references companies
            on delete restrict,
    avatar_url  text,
    created_at  timestamp with time zone,
    updated_at  timestamp with time zone
);

alter table companies
    add constraint fk_companies_owner
        foreign key (owner_id) references users
            on delete restrict;

create index idx_users_id
    on users (id);


COMMIT;