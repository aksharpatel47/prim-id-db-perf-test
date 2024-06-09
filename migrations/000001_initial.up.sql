BEGIN;

create extension if not exists "uuid-ossp";

CREATE TABLE IF NOT EXISTS person_bigserial
(
    id         bigserial PRIMARY KEY NOT NULL,
    first_name varchar(255)          NOT NULL,
    last_name  varchar(255)          NOT NULL,
    data_1     varchar(255)          NOT NULL,
    data_2     varchar(255)          NOT NULL,
    data_3     varchar(255)          NOT NULL,
    data_4     varchar(255)          NOT NULL,
    data_5     varchar(255)          NOT NULL
);

CREATE TABLE IF NOT EXISTS person_random_int
(
    id         bigint PRIMARY KEY NOT NULL,
    first_name varchar(255)       NOT NULL,
    last_name  varchar(255)       NOT NULL,
    data_1     varchar(255)       NOT NULL,
    data_2     varchar(255)       NOT NULL,
    data_3     varchar(255)       NOT NULL,
    data_4     varchar(255)       NOT NULL,
    data_5     varchar(255)       NOT NULL
);

CREATE TABLE IF NOT EXISTS person_date_random_int
(
    id         bigint PRIMARY KEY NOT NULL,
    first_name varchar(255)       NOT NULL,
    last_name  varchar(255)       NOT NULL,
    data_1     varchar(255)       NOT NULL,
    data_2     varchar(255)       NOT NULL,
    data_3     varchar(255)       NOT NULL,
    data_4     varchar(255)       NOT NULL,
    data_5     varchar(255)       NOT NULL
);

CREATE table if not exists person_uuid
(
    id         uuid PRIMARY KEY NOT NULL default gen_random_uuid(),
    first_name varchar(255)     NOT NULL,
    last_name  varchar(255)     NOT NULL,
    data_1     varchar(255)     NOT NULL,
    data_2     varchar(255)     NOT NULL,
    data_3     varchar(255)     NOT NULL,
    data_4     varchar(255)     NOT NULL,
    data_5     varchar(255)     NOT NULL
);

CREATE table if not exists person_uuidv7
(
    id         varchar(36) primary key not null,
    first_name varchar(255)            NOT NULL,
    last_name  varchar(255)            NOT NULL,
    data_1     varchar(255)            NOT NULL,
    data_2     varchar(255)            NOT NULL,
    data_3     varchar(255)            NOT NULL,
    data_4     varchar(255)            NOT NULL,
    data_5     varchar(255)            NOT NULL
);

CREATE table if not exists addresses_bigserial
(
    id        bigserial primary key                   not null,
    person_id bigint references person_bigserial (id) not null,
    address   varchar(255)                            not null,
    city      varchar(255)                            not null,
    state     varchar(255)                            not null,
    zip       varchar(255)                            not null
);

CREATE table if not exists addresses_random_int
(
    id        bigserial primary key                    not null,
    person_id bigint references person_random_int (id) not null,
    address   varchar(255)                             not null,
    city      varchar(255)                             not null,
    state     varchar(255)                             not null,
    zip       varchar(255)                             not null
);

CREATE table if not exists addresses_date_random_int
(
    id        bigserial primary key                         not null,
    person_id bigint references person_date_random_int (id) not null,
    address   varchar(255)                                  not null,
    city      varchar(255)                                  not null,
    state     varchar(255)                                  not null,
    zip       varchar(255)                                  not null
);

CREATE table if not exists addresses_uuid
(
    id        bigserial primary key            not null,
    person_id uuid references person_uuid (id) not null,
    address   varchar(255)                     not null,
    city      varchar(255)                     not null,
    state     varchar(255)                     not null,
    zip       varchar(255)                     not null
);

CREATE table if not exists addresses_uuidv7
(
    id        bigserial primary key                     not null,
    person_id varchar(36) references person_uuidv7 (id) not null,
    address   varchar(255)                              not null,
    city      varchar(255)                              not null,
    state     varchar(255)                              not null,
    zip       varchar(255)                              not null
);

COMMIT;