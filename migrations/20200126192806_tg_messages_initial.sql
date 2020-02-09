-- +goose Up
-- +goose StatementBegin
create table tg_messages
(
    created    integer not null,
    message_id integer not null,
    offer_id   integer not null,
    chat       integer not null
);

create index tg_messages_message_id_index
    on tg_messages (message_id);

create index tg_messages_offer_id_index
    on tg_messages (offer_id);

create index tg_messages_chat_index
    on tg_messages (chat);

create unique index tg_messages_message_id_offer_id_chat_uindex
    on tg_messages (message_id, offer_id, chat);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists tg_messages;
-- +goose StatementEnd
