CREATE TYPE transaction_type as ENUM ('c', 'd');

CREATE TABLE transactions
(
    id          BIGSERIAL PRIMARY KEY,
    customer_id BIGINT           NOT NULL,
    amount      SERIAL           NOT NULL,
    description VARCHAR(10)      NOT NULL,
    type        transaction_type NOT NULL,
    created_at  timestamptz      NOT NULL,

    CONSTRAINT fk_customer FOREIGN KEY (customer_id) REFERENCES customers (id)
);
