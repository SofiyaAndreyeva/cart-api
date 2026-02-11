-- +goose Up
-- +goose StatementBegin
CREATE TABLE carts
(
    id         SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE cart_items
(
    id      BIGSERIAL PRIMARY KEY,
    cart_id INT REFERENCES carts (id) ON DELETE CASCADE NOT NULL,
    product TEXT                                        NOT NULL CHECK (product <> ''),
    price   NUMERIC(10, 2)                              NOT NULL CHECK (price > 0)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cart_items;

DROP TABLE IF EXISTS carts;
-- +goose StatementEnd