-- +goose Up
create table offers (
id serial primary key,
offer_id bigint unique,
name varchar(25),
price float,
available bool
);

-- +goose Down
drop table offers;
