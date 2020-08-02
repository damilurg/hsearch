-- +goose Up
-- +goose StatementBegin
alter table chat
    add house integer default 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
create table chat_dg_tmp
(
    id integer not null,
    username varchar(100) default '',
    enable integer default 0,
    title varchar(100) default '',
    c_type varchar(20) default '',
    created integer default 0,
    diesel integer default 1,
    lalafo integer default 1,
    photo integer default 0,
    usd varchar(100) default '0:0',
    kgs varchar(100) default '0:0',
    up_track integer default 0
);

insert into chat_dg_tmp(id, username, enable, title, c_type, created, diesel, lalafo, photo, usd, kgs, up_track) select id, username, enable, title, c_type, created, diesel, lalafo, photo, usd, kgs, up_track from chat;

drop table chat;

alter table chat_dg_tmp rename to chat;

create unique index chat_chat_index
    on chat (id);
-- +goose StatementEnd
