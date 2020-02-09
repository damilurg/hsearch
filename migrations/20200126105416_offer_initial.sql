-- +goose Up
-- +goose StatementBegin
create table offer
(
    id           integer      not null,
    created      integer      not null,
    url          varchar(255) default '',
    topic        varchar(255) default '',
    price        real,
    phone        varchar(255) default '',
    room_numbers varchar(255) default '',
    body         text         default '',
    images       integer      not null
);

create unique index offer_id_uindex
    on offer (id);


-- sqlite не умеет хранить списки в полях, по этому делам отдельную таблицу
create table image
(
    offer_id integer      not null,
    path     varchar(255) default ''
);

create index image_offer_id_index
    on image (offer_id);

create unique index image_path_uindex
	on image (path);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists offer;
drop table if exists image;
-- +goose StatementEnd
