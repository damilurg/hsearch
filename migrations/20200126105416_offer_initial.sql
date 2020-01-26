-- +goose Up
-- +goose StatementBegin
create table offer
(
    id    integer not null
        constraint offer_pk
            primary key autoincrement,
    topic varchar(255) default '',
    body  text         default '',
    price real,
    ex_id integer not null
);

create unique index offer_ex_id_uindex
    on offer (ex_id);

create unique index offer_id_uindex
    on offer (id);


-- sqlite не умеет хранить списки в полях, по этому делам отдельную таблицу
create table image
(
    offer_id integer not null,
    path     varchar(255) default ''
);

create index image_offer_id_index
    on image (offer_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists offer;
drop table if exists image;
-- +goose StatementEnd
