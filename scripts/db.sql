CREATE DATABASE calendar;

CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS customer
(
    id    uuid primary key,
    email citext not null unique
);

CREATE TABLE IF NOT EXISTS event
(
    id                  uuid primary key,
    owner_id            uuid references customer (id) not null,
    organizer_full_name text                          not null,
    organizer_email     citext                        not null,
    time_start          text                          not null,
    time_end            text                          not null,
    event_summary       text                          not null,
    geo_lat             real                          not null,
    geo_lon             real                          not null
);
