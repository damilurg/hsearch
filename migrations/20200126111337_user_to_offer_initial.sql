-- +goose Up
-- +goose StatementBegin
create table user_to_offer
(
    user_id  integer not null,
    offer_id integer not null,
    -- sqlite не поддерживает boolean тип, по этому спользуется integer 1/0
    -- в других бд, boolean выглядит как tinyint, то есть 1/0, так что это норм
    like     integer default 0,
    dislike  integer default 0,
    skip     integer
);

create index user_to_offer_offer_id_index
    on user_to_offer (offer_id);

create index user_to_offer_user_id_index
    on user_to_offer (user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists user_to_offer;
-- +goose StatementEnd
