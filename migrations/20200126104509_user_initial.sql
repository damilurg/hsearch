-- +goose Up
-- +goose StatementBegin
create table user
(
    id       integer not null
        constraint user_pk
            primary key autoincrement,
    username varchar(100),
    chat     integer,
    enable   integer default 0
);


create unique index user_id_uindex
    on user (id);

create index user_chat_index
    on user (chat);

create unique index user_name_uindex
    on user (username);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists user;
-- +goose StatementEnd
