CREATE TABLE stocks
(
    id serial not null unique,
    name varchar(255) not null,
    feature varchar(255) not null default 'available'
);

CREATE TABLE products
(
    id serial not null unique,
    name varchar(255) not null,
    size integer not null default 0,
    uniq_code integer not null default 0,
    count integer not null default 0
);

CREATE TABLE stocks_products
(
    id serial not null unique,
    stocks_id int references stocks(id) on delete cascade not null,
    products_id int references products(id) on delete cascade not null
);