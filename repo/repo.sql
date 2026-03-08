DROP TABLE IF EXISTS users, store, activity, track_totals;

CREATE TABLE users (
    id serial primary key,
    nkey varchar(255),
    user_name varchar(255)
);

CREATE TABLE store (
    id serial primary key,
    business_name varchar(255),
    point_step integer,
    point_to_cash decimal
);

CREATE TABLE activity (
    user_id integer not null,
    foreign key(user_id) references users(id) on delete cascade,
    store_id integer not null,
    foreign key(store_id) references store(id) on delete cascade,
    visit_date date, /* calc (date added of store) */
    point_redeem integer,
    cash_redeem decimal,
    points_added integer
);

CREATE TABLE track_totals (
    user_id integer not null,
    foreign key(user_id) references users(id) on delete cascade,
    store_id integer not null,
    foreign key(store_id) references store(id) on delete cascade,
    total_points integer,
    total_cash decimal
);