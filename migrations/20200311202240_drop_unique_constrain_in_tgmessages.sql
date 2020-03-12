-- +goose Up
-- +goose StatementBegin
drop index tg_messages_message_id_offer_id_chat_uindex;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
create unique index tg_messages_message_id_offer_id_chat_uindex
    on tg_messages (message_id, offer_id, chat);
-- +goose StatementEnd
