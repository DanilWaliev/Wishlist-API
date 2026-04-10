BEGIN;

create table users
(
	id serial primary key,
	email varchar(255) not null unique,
	password_hash varchar(150) not null
);

create table wishlists
(
	id serial primary key,
	event_name varchar(50) not null,
	description text,
	event_date date not null,
	token text not null unique,
	user_id integer not null references users(id) on delete cascade
);

create table wishlist_items
(
	id serial primary key,
	title varchar(50) not null,
	description text,
	product_url text,
	priority integer not null default 1 check (priority >= 1 and priority <= 10),
	reserved boolean not null default false,
	wishlist_id integer not null references wishlists(id) on delete cascade
);

COMMIT;