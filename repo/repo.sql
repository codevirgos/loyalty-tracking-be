
CREATE DATABASE PB;

DROP TABLE users;

CREATE TABLE users (
    nkey integer,
    last_visit date,
    points integer,
    visit_count integer,
    name varchar(255)
);

CREATE TABLE store {
    store id integer,
    point_step integer,
}