-- +goose Up
-- +goose StatementBegin
create table save_message
(
    message_id integer not null,
    offer_id   integer not null,
    chat       integer not null
);

create index save_message_offer_id_index
    on save_message (message_id);

create index save_message_user_id_index
    on save_message (offer_id);

create index save_message_chat_index
    on save_message (chat);

create unique index save_message_message_id_user_id_chat_uindex
    on save_message (message_id, offer_id, chat);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists save_message;
-- +goose StatementEnd
