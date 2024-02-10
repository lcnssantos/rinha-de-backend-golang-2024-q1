CREATE TABLE customers
(
    id             BIGSERIAL PRIMARY KEY,
    "limit"        INT NOT NULL,
    amount INT NOT NULL,
    CONSTRAINT customer_limit_check CHECK (-amount <= "limit")
);

ALTER DATABASE postgres SET log_error_verbosity to 'TERSE';

