-- +goose Up
-- +goose StatementBegin
create table answer_dg_tmp
(
    created integer not null,
    chat integer not null,
    offer_id integer not null,
    like integer default 0,
    dislike integer default 0
);

insert into answer_dg_tmp(created, chat, offer_id, like, dislike) select created, chat, offer_id, like, dislike from answer;

drop table answer;

alter table answer_dg_tmp rename to answer;

create index answer_chat_index
    on answer (chat);

create index answer_offer_id_index
    on answer (offer_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table answer
    add skip integer;

-- +goose StatementEnd
