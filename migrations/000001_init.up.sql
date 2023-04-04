CREATE TABLE stocks
(
    id serial not null unique,
    name varchar(255) not null,
    feature varchar(255) not null
);

CREATE TABLE products
(
    id serial not null unique,
    name varchar(255) not null,
    size integer not null,
    uniq_code integer not null,
    count integer not null
);

CREATE TABLE stocks_products
(
    id serial not null unique,
    stocks_id int references stocks(id) on delete cascade not null,
    products_id int references products(id) on delete cascade not null
);