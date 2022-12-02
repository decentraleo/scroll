-- +goose Up
-- +goose StatementBegin
create table l2_message
(
    nonce        BIGINT  NOT NULL,
    height       BIGINT  NOT NULL,
    sender       VARCHAR NOT NULL,
    target       VARCHAR NOT NULL,
    value        VARCHAR NOT NULL,
    fee          VARCHAR NOT NULL,
    gas_limit    BIGINT  NOT NULL,
    deadline     BIGINT  NOT NULL,
    calldata     TEXT    NOT NULL,
    layer2_hash  VARCHAR NOT NULL,
    layer1_hash  VARCHAR DEFAULT NULL,
    proof        TEXT    DEFAULT NULL,
    status       INTEGER  DEFAULT 1,
    created_time TIMESTAMP(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_time TIMESTAMP(0) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

comment
on column l2_message.status is 'undefined, pending, submitted, confirmed';

create unique index l2_message_layer2_hash_uindex
    on l2_message (layer2_hash);

create index l2_message_height_index
    on l2_message (height);

CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_time = CURRENT_TIMESTAMP;
   RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_timestamp BEFORE UPDATE
ON l2_message FOR EACH ROW EXECUTE PROCEDURE
update_timestamp();


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists l2_message;
-- +goose StatementEnd