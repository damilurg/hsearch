-- +goose Up
-- +goose StatementBegin
create table user
(
    id      integer not null
        constraint user_pk
            primary key autoincrement,
    account varchar(100),
    chat    int
);


create unique index user_id_uindex
    on user (id);

create index user_chat_index
    on user (chat);

create unique index user_name_uindex
    on user (account);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists user;
-- +goose StatementEnd
