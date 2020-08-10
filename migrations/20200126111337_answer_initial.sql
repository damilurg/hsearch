-- +goose Up
-- +goose StatementBegin
create table answer
(
    created  integer not null,
    chat     integer not null,
    offer_id integer not null,
    -- sqlite не поддерживает boolean тип, по этому используется integer 1/0
    -- в других бд, boolean выглядит как tinyint, то есть 1/0, так что это норм
    like     integer default 0,
    dislike  integer default 0,
    skip     integer
);

create index answer_offer_id_index
    on answer (offer_id);

create index answer_chat_index
    on answer (chat);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists answer;
-- +goose StatementEnd
